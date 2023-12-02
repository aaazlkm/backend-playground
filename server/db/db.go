// package db
package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

type DB struct {
	db *sql.DB
}

func NewDB() (*DB, error) {
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

		return nil, fmt.Errorf("database ping error %w", err)
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

func (d *DB) Create(album Album) (int, error) {
	result, err := d.db.Exec(
		"insert into album (title, artist, price) values (?, ?, ?)",
		album.Title,
		album.Artist,
		album.Price,
	)
	if err != nil {
		return 0, fmt.Errorf("database create error %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("database create error %w", err)
	}

	return int(id), err
}

func (d *DB) Update(album Album) (int, error) {
	result, err := d.db.Exec(
		"update album set title = ?, artist = ?, price = ? WHERE id = ?",
		album.Title,
		album.Artist,
		album.Price,
		album.ID,
	)
	if err != nil {
		return 0, fmt.Errorf("database update error %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("database update error %w", err)
	}

	return int(id), err
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

func (d *DB) ReadAll() ([]Album, error) {
	var albums []Album

	rows, err := d.db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("database read error %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("database read error %w", err)
		}

		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("database read error %w", err)
	}

	return albums, nil
}

/* -------------------------------------------------------------------------- */
/*                                    model                                   */
/* -------------------------------------------------------------------------- */

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}
