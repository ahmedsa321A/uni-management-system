package store

import (
	"database/sql"
	"fmt"
	"university-management/backend/models"
)

type BOOKstore struct {
	DB *sql.DB
}

func (s *BOOKstore) Create(book *models.Book) (int64, error) {
	query := `
		INSERT INTO BOOKS (title, author, published_year, isbn, publisher)
		VALUES (?, ?, ?, ?);`
	result, err := s.DB.Exec(query, book.Title, book.Author, book.PublicationYear, book.ISBN, book.Publisher)
	if err != nil {
		return 0, err
	}
	newBookID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return newBookID, nil
}
func (s *BOOKstore) GetByID(bookID int) (*models.Book, error) {
	query := `
		SELECT 
			book_id, title, author, published_year, isbn, publisher	
		WHERE
			book_id = ?;`

	row := s.DB.QueryRow(query, bookID)
	var book models.Book
	err := row.Scan(
		&book.BookID,
		&book.Title,
		&book.Author,
		&book.PublicationYear,
		&book.ISBN,
		&book.Publisher,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no book found with ID %d", bookID)
		}
		return nil, err
	}
	return &book, nil
}

// get books
func (s *BOOKstore) GetAll() ([]*models.Book, error) {
	query := `
		SELECT 
			book_id, title, author, published_year, isbn, publisher
		FROM 
			BOOKS;`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []*models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.BookID,
			&book.Title,
			&book.Author,
			&book.PublicationYear,
			&book.ISBN,
			&book.Publisher,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}
