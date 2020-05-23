package db

import (
	"database/sql"

	"github.com/ellemouton/snell/articles"
)

func Create(dbc *sql.DB, name, description string, price int64, text string) (int64, error) {
	// TODO(elle): Replace with a transaction.These two inserts must be atomic
	res, err := dbc.Exec("insert into articles_content set text=?", text)
	if err != nil {
		return 0, err
	}

	contentID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	res, err = dbc.Exec("insert into articles_info set name=?, description=?, "+
		"price=?, content_id=?", name, description, price, contentID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func LookupInfo(dbc *sql.DB, id int64) (*articles.Info, error) {
	row := dbc.QueryRow("select * from articles_info where id=?", id)

	info := articles.Info{}
	err := row.Scan(&info.ID, &info.Name, &info.Description, &info.Price, &info.ContentID)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func LookupContent(dbc *sql.DB, id int64) (*articles.Content, error) {
	row := dbc.QueryRow("select * from articles_content where id=?", id)

	content := articles.Content{}
	err := row.Scan(&content.ID, &content.Text)
	if err != nil {
		return nil, err
	}

	return &content, nil
}

func ListAllInfo(dbc *sql.DB) (infos []*articles.Info, err error) {
	rows, err := dbc.Query("select * from articles_info")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		info := articles.Info{}
		err = rows.Scan(&info.ID, &info.Name, &info.Description, &info.Price, &info.ContentID)
		if err != nil {
			return nil, err
		}
		infos = append(infos, &info)
	}
	return infos, rows.Err()
}
