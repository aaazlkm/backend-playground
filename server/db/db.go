// package db
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type DB struct {
	db *sql.DB
}

func NewDB(count int) (*DB, error) {
	if count == 0 {
		return nil, fmt.Errorf("database connect error, retry count is 0")
	}

	path := fmt.Sprintf(
		"%s:%s@tcp(mysql:3306)/%s?charset=utf8&parseTime=true",
		`user`,
		`password`,
		`db`,
	)

	db, err := sql.Open("mysql", path)
	if err != nil {
		slog.Error(`database open error %w`, err)

		return nil, fmt.Errorf("database Open error %w", err)
	}

	if err := db.Ping(); err != nil {
		slog.Error(`database connect error %v`, err)

		time.Sleep(time.Second * 2)
		count--
		slog.Info("retry... count:%d\n", count)

		return NewDB(count)
	}

	slog.Info("Connected to DB")

	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("database close error %w", err)
	}

	return nil
}

/* -------------------------------------------------------------------------- */
/*                                    crud                                   */
/* -------------------------------------------------------------------------- */

func (d *DB) Create(album *Album) error {
	// var alb Album

	// row := d.db.QueryRow("SELECT * FROM album WHERE id = ?")

	return nil
}

func (d *DB) Read(id int) (Album, error) {
	var alb Album

	row := d.db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return alb, fmt.Errorf("no album with id %d", id)
		}

		return alb, fmt.Errorf("database read error %w", err)
	}

	return alb, nil
}

/* -------------------------------------------------------------------------- */
/*                                    model                                   */
/* -------------------------------------------------------------------------- */
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}
