package postgres

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"grpc/pkg/user/model"
	"grpc/pkg/user/storage"
	"log"
)

type Postgres struct {
	db *sql.DB
}

func (p Postgres) CreateUser(ctx context.Context, user *model.UserDB) (int, error) {

	stmt := `INSERT INTO users (email) values $1 RETURNING id`
	var id int
	err := p.db.QueryRow(stmt, user.Email).Scan(id)
	if err != nil {
		return 0, err
	}
	return id, nil

}

func (p Postgres) GetAll(ctx context.Context, offset, limit int) ([]*model.UserDB, error) {
	stmt := `SELECT (id,email) FROM users`
	rows, err := p.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	var users []*model.UserDB
	for rows.Next() {
		var user *model.UserDB
		err := rows.Scan(&user.ID, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (p Postgres) DeleteUser(ctx context.Context, id int) error {
	_, err := p.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

var _ storage.CRUDL = &Postgres{}

func NewPostgres(connStr string) (*Postgres, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	m, err := migrate.New(
		"file://db/migration/postgresql",
		connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal("error while upp %s   ", err.Error())
	}

	return &Postgres{db: db}, nil

}
