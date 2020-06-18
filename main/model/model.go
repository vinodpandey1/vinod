// model.go

package model

import (
	"database/sql"
	"fmt"
)

//Book Json request payload is as follows,
//{
//  "id": "1",
//  "author": "james",
//  "title":  "bolt",
//  "isbn":  "james1234"
//  "release_date": "2016-10-10"
//}
type Book struct {
	ID          int    `json:"id"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	ISBN        string `json:"isbn"`
	ReleaseDate string `json:"release_date"`
}

// Book response payload
// swagger:response bookRes
type swaggBookRes struct {
	// in:body
	Body Book
}

// Error Bad Request
// swagger:response badReq
type swaggReqBadRequest struct {
	// in:body
	Body struct {
		// HTTP status code 400 -  Bad Request
		Code int `json:"code"`
	}
}

// Error Not Found
// swagger:response notFoundReq
type swaggReqNotFound struct {
	// in:body
	Body struct {
		// HTTP status code 404 -  Not Found
		Code int `json:"code"`
	}
}

func (u *Book) GetBook(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT author, title,isbn,release_date FROM books WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Author, &u.Title, &u.ISBN, &u.ReleaseDate)
}

func (u *Book) UpdateBook(db *sql.DB) error {
	//*u.ReleaseDate = time.Now()
	statement := fmt.Sprintf("UPDATE books SET author='%s', title='%s', isbn='%s',release_date='%s' WHERE id=%d", u.Author, u.Title, u.ISBN, u.ReleaseDate, u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *Book) DeleteBook(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM books WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *Book) CreateBook(db *sql.DB) error {
	//u.ReleaseDate = *time.Now()
	statement := fmt.Sprintf("INSERT INTO books(author, title,isbn,release_date ) VALUES('%s', '%s','%s','%s')", u.Author, u.Title, u.ISBN, u.ReleaseDate)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetBooks(db *sql.DB, start, count int) ([]Book, error) {
	statement := fmt.Sprintf("SELECT id, author, title,isbn,release_date FROM books LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []Book{}

	for rows.Next() {
		var u Book
		if err := rows.Scan(&u.ID, &u.Author, &u.Title, &u.ISBN, &u.ReleaseDate); err != nil {
			return nil, err
		}
		books = append(books, u)
	}

	return books, nil
}
