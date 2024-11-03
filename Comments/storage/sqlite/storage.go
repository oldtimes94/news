package sqlite

import (
	"comments/models"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Storage struct {
	DB *sql.DB
}

func New() *Storage {
	db, err := sql.Open("sqlite3", "./comments.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		news_id INTEGER,
		parent_comment_id INTEGER,
		text TEXT,
		author TEXT
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	return &Storage{DB: db}

}

func (s *Storage) Add(newsID, parentID int, text string) error {

	stmt, err := s.DB.Prepare("INSERT INTO comments(news_id, parent_comment_id, text) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(newsID, parentID, text)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CommentByNewsID(newsID int) ([]models.Comment, error) {
	rows, err := s.DB.Query("SELECT id, text FROM comments WHERE news_id=?", newsID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := make([]models.Comment, 0, 10)
	for rows.Next() {
		comment := models.Comment{}
		err = rows.Scan(&comment.Id, &comment.Text)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}
	return comments, nil
}
