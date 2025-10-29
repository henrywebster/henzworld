// Package database: handling database
package database

import (
	"bytes"
	"database/sql"
	"henzworld/internal/model"
	"html/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
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

func (db *DB) GetPost(slug string) (*model.Post, error) {
	query := "SELECT title, DATE(created_at) as created_at, content FROM posts WHERE slug = ?"

	var post model.Post
	var createdAtString string
	var contentReadBuffer []byte
	row := db.conn.QueryRow(query, slug)
	if err := row.Scan(&post.Title, &createdAtString, &contentReadBuffer); err != nil {
		return nil, err
	}
	formattedCreatedAt, err := time.Parse("2006-01-02", createdAtString)
	if err != nil {
		return nil, err
	}
	post.CreatedAt = formattedCreatedAt

	markdown := goldmark.New(goldmark.WithExtensions(extension.Footnote))

	var contentWriteBuffer bytes.Buffer
	if err := markdown.Convert(contentReadBuffer, &contentWriteBuffer); err != nil {
		return nil, err
	}
	post.Content = template.HTML(contentWriteBuffer.String())

	return &post, nil
}

func (db *DB) GetPosts() ([]model.Post, error) {
	query := "SELECT title, created_at, description, slug FROM posts ORDER BY created_at DESC"

	var posts []model.Post
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post
		var createdAtString string
		if err := rows.Scan(&post.Title, &createdAtString, &post.Description, &post.Slug); err != nil {
			return nil, err
		}
		formattedCreatedAt, err := time.Parse("2006-01-02 15:04:05", createdAtString)
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

func (db *DB) GetText(name string) (*template.HTML, error) {
	query := "SELECT content FROM texts WHERE name = ?"
	var contentReadBuffer []byte

	row := db.conn.QueryRow(query, name)
	if err := row.Scan(&contentReadBuffer); err != nil {
		return nil, err
	}

	markdown := goldmark.New(goldmark.WithExtensions(extension.Footnote))

	var contentWriteBuffer bytes.Buffer
	if err := markdown.Convert(contentReadBuffer, &contentWriteBuffer); err != nil {
		return nil, err
	}
	content := template.HTML(contentWriteBuffer.String())
	return &content, nil
}
