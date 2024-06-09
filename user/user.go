package user

import (
	"database/sql"
	"time"
)

type user struct {
	Id        int64
	Username  string
	Password  string
	CreatedAt time.Time
}

type UserRepository interface {
	CreateUser(username string, password string) (*user, error)
	GetUserById(id int64) (*user, error)
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

func (repo sqliteUserRepository) CreateUser(username string, password string) (*user, error) {
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
		Id:        newId,
		Username:  username,
		Password:  password,
		CreatedAt: createdAt,
	}
	return new_user, nil
}

func (repo sqliteUserRepository) GetUserById(id int64) (*user, error) {
	user := user{}
	query := "SELECT rowid, username, password, created_at FROM users WHERE rowid = ?"
	err := repo.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
