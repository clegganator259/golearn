package user

import (
	"database/sql"
	"time"
)

type user struct {
	id        int64
	username  string
	password  string
	createdAt time.Time
}

type UserRepository interface {
	createUser(username string, password string) (*user, error)
	getUserById(id int64) (*user, error)
}

type sqliteUserRepository struct {
	db *sql.DB
}

func NewSqliteRepo(connection_string string) (UserRepository, error) {
	db, err := sql.Open("sqlite3", connection_string)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return &sqliteUserRepository{db}, nil
}

func (repo sqliteUserRepository) createUser(username string, password string) (*user, error) {
	createdAt := time.Now()
	result, err := repo.db.Exec(`INSERT INTO users (username, password, created_at) values (?, ?, ?)`, username, password, createdAt)
	if err != nil {
		return nil, err
	}
	newId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	new_user := &user{
		id:        newId,
		username:  username,
		password:  password,
		createdAt: createdAt,
	}
	return new_user, nil
}

func (repo sqliteUserRepository) getUserById(id int64) (*user, error) {
	user := user{}
	query := "SELECT rowid, username, password, created_at FROM users WHERE rowid = ?"
	err := repo.db.QueryRow(query, id).Scan(&user.id, &user.username, &user.password, &user.createdAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
