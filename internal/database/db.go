// Package database: handling database
package database

import (
	"database/sql"
	"henzworld/internal/model"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	conn *sql.DB
}

func New(dbSourceName string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dbSourceName)
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) GetPosts() ([]model.Post, error) {
	query := "SELECT title, created_at, description FROM posts ORDER BY created_at DESC"

	var posts []model.Post
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		var createdAtString string
		if err := rows.Scan(&post.Title, &createdAtString, &post.Description); err != nil {
			return nil, err
		}
		formattedCreatedAt, err := time.Parse("2006-01-02 15:04", createdAtString)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = formattedCreatedAt

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
