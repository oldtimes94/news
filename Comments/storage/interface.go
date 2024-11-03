package storage

import "comments/models"

type Interface interface {
	Add(newsID, parentID int, text string) error
	CommentByNewsID(newsID int) ([]models.Comment, error)
}
