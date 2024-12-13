package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

type User struct {
	Id       int64
	ChatId   int64
	Username string
	Root     bool
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Init() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id BIGINT PRIMARY KEY,
		chat_id BIGINT,
		username TEXT,
		root BOOL
	);`

	_, err := s.db.Exec(query)
	return err
}

func (s *Store) SaveUser(user User) error {
	query := `INSERT INTO users (id, chat_id, username,  root) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET chat_id=$2, username=$3, root=$4;`
	_, err := s.db.Exec(query, user.Id, user.ChatId, user.Username, user.Root)
	return err
}

func (s *Store) GetUsers() ([]User, error) {
	query := `SELECT id, chat_id, username, root FROM users;`
	users := make([]User, 0, 1000)
	rows, err := s.db.Query(query)
	for rows.Next() {
		user := User{}
		rows.Scan(&user)
		users = append(users, user)
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}
