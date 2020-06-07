package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/skip2/go-qrcode"

	articles_db "github.com/ellemouton/snell/articles/db"
)

func registerRoutes(s *State) {
	router.GET("/", s.homeHandler)
	router.GET("/post", s.postArticleFormHandler)
	router.GET("/article/view/:id", s.viewArticleHandler)
	router.POST("/post", s.saveArticleHandler)
}

func (s *State) homeHandler(c *gin.Context) {
	articles, err := articles_db.ListAllInfo(s.GetDB())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":   "",
			"payload": articles,
		},
	)
}

func (s *State) postArticleFormHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"post.html",
		gin.H{
			"title": "-post",
		},
	)
}

func (s *State) saveArticleHandler(c *gin.Context) {
	title := c.PostForm("title")
	abstract := c.PostForm("abstract")
	content := c.PostForm("content")
	fmt.Println(c.PostForm("price"))
	price, err := strconv.ParseInt(c.PostForm("price"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	// TODO(elle): Create ops package for articles rather than calling db layer
	_, err = articles_db.Create(s.GetDB(), title, abstract, price, content)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	s.homeHandler(c)
}

var authRegex = regexp.MustCompile("LSAT (.*?):([a-f0-9]{64})")

func (s *State) viewArticleHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	article, err := articles_db.LookupInfo(s.GetDB(), id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if article.Price == 0 {
		s.displayArticleHandler(c)
		return
	}

	// Check header for authorization and redirect to payment handler if needed.
	auth := c.GetHeader("Authorization")

	if auth == "" {
		s.paymentHandler(c)
		return
	}

	if !authRegex.MatchString(auth) {
		s.paymentHandler(c)
		return
	}

	matches := authRegex.FindStringSubmatch(auth)
	if len(matches) != 3 {
		s.paymentHandler(c)
		return
	}

	macString, preimageHexString := matches[1], matches[2]

	macBytes, err := base64.StdEncoding.DecodeString(macString)
	if err != nil {
		s.paymentHandler(c)
		return
	}

	preimage, err := hex.DecodeString(preimageHexString)
	if err != nil {
		s.paymentHandler(c)
		return
	}

	valid, err := s.macClient.Verify(macBytes, preimage, "article", id)
	if err != nil {
		s.paymentHandler(c)
		return
	}

	if !valid {
		s.paymentHandler(c)
		return
	}

	s.displayArticleHandler(c)
}

func (s *State) displayArticleHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	article, err := articles_db.LookupInfo(s.GetDB(), id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	content, err := articles_db.LookupContent(s.GetDB(), article.ContentID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.HTML(
		http.StatusOK,
		"article.html",
		gin.H{
			"title":   "-view",
			"info":    article,
			"payload": content,
		},
	)
}

func (s *State) paymentHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	article, err := articles_db.LookupInfo(s.GetDB(), id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	// Generate a new invoice to be paid
	invoice, err := s.lndClient.AddInvoice(context.Background(), article.Price, 3600, article.Name)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var payHash lntypes.Hash
	copy(payHash[:], invoice.RHash)

	// Bake a new macaroon
	mac, err := s.macClient.Create(payHash, "article", article.ID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	macBytes, err := mac.MarshalBinary()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	macString := base64.StdEncoding.EncodeToString(macBytes)

	// construct QR code of the invoice
	png, err := qrcode.Encode(strings.ToUpper(invoice.PaymentRequest), qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

	encodedPngString := base64.StdEncoding.EncodeToString(png)

	// Add the partial LSAT (mac + invoice) to the response header
	str := fmt.Sprintf("LSAT macaroon=\"%s\", invoice=\"%s\"", macString, invoice.PaymentRequest)
	c.Writer.Header().Set("WWW-Authenticate", str)

	c.HTML(
		http.StatusPaymentRequired,
		"payment.html",
		gin.H{
			"title":    "-pay",
			"article":  article,
			"invoice":  invoice.PaymentRequest,
			"macaroon": macString,
			"qrCode":   encodedPngString,
		},
	)
}

func (s *State) checkForAuthHandler(c *gin.Context) {
	auth := c.GetHeader("Authorization")

	if auth == "" {
		s.paymentHandler(c)
		return
	}

	if !authRegex.MatchString(auth) {
		s.paymentHandler(c)
		return
	}

	matches := authRegex.FindStringSubmatch(auth)
	if len(matches) != 3 {
		s.paymentHandler(c)
		return
	}

	macBytes, preimageHex := matches[1], matches[2]
	fmt.Fprintln(c.Writer, macBytes, preimageHex)
}
