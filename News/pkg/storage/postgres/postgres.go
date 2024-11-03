package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

// Хранилище данных.
type Store struct {
	Pool *pgxpool.Pool
}

func New(conStr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), conStr)
	if err != nil {
		return nil, err
	}
	s := Store{
		Pool: db,
	}
	return &s, nil
}

type NewsPost struct {
	ID      int    // номер записи
	Title   string // заголовок публикации
	Content string // содержание публикации
	PubTime int64  // время публикации
	Link    string // ссылка на источник
}

func (s *Store) Posts(query string) ([]storage.NewsPost, error) {
	posts := make([]storage.NewsPost, 0, 50)
	rows, err := s.Pool.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post storage.NewsPost
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *Store) Post(id int) (storage.NewsPost, error) {
	rows, err := s.Pool.Query(
		context.Background(),
		`SELECT news.id, title, content, pubtime, link FROM news WHERE id = $1`,
		id,
	)
	if err != nil {
		return storage.NewsPost{}, err
	}

	var post storage.NewsPost
	for rows.Next() {
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.PubTime,
			&post.Link,
		)
		if err != nil {
			return storage.NewsPost{}, err
		}
	}

	return post, nil
}

func (s *Store) Count(query string) (int, error) {
	var totalCount int
	err := s.Pool.QueryRow(
		context.Background(),
		query,
	).Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil

}

func (s *Store) AddPosts(posts []storage.NewsPost) error {

	ctx := context.Background()

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)

	for _, post := range posts {
		batch.Queue(`INSERT INTO news(title, content, pubtime, link) VALUES 
                                                    ($1, $2, $3, $4)`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link)
	}

	res := tx.SendBatch(ctx, batch)
	// обязательная операция закрытия соединения
	err = res.Close()
	if err != nil {
		return err
	}
	// подтверждение транзакции
	return tx.Commit(ctx)

}
