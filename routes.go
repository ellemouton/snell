package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lightningnetwork/lnd/lntypes"

	articles_db "github.com/ellemouton/snell/articles/db"
)

func registerRoutes(s *State) {
	router.GET("/", s.homeHandler)
	router.GET("/post", s.postArticleFormHandler)
	router.GET("/article/view/:id", s.viewArticleHandler)

	router.GET("/fancy", s.paymentHandler)
	router.GET("/blah", s.checkForAuthHandler)

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
			"title":   "Home",
			"payload": articles,
		},
	)
}

func (s *State) postArticleFormHandler(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"post.html",
		gin.H{
			"title": "Post Article",
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

	// macBytes, preimageHex := matches[1], matches[2]

	// validate:
	//	1. check mac was minted by me
	//	2. validate preimage
	//	3. check caveat for given article

	// not valid: payment handler
	// else: view page

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
			"title":   "View",
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

	str := fmt.Sprintf("LSAT macaroon=\"%s\", invoice=\"%s\"", base64.StdEncoding.EncodeToString(macBytes), invoice.PaymentRequest)
	c.Writer.Header().Set("WWW-Authenticate", str)

	c.HTML(
		http.StatusPaymentRequired,
		"payment.html",
		gin.H{
			"title": "payment",
			"LSAT":  str,
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
