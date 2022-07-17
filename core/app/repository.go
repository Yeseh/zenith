package app

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
	zd "github.com/yeseh/zenith/domain"
)

func (r AppRepository) New(db *sql.DB) *AppRepository {
	return &AppRepository{
		db,
	}
}

func (r *AppRepository) Migrate() error {
	query := `
		CREATE TABLE IF NOT EXISTS deployments(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			displayName TEXT NOT NULL UNIQUE,
			location TEXT NOT NULL,
			runtime TEXT NOT NULL,
			images TEXT
		);`

	_, err := r.db.Exec(query)

	return err
}

func (r *AppRepository) GetByDisplayName(name string) (zd.App, error) {
	query := `SELECT * FROM deployments WHERE displayName = ?`
	row := r.db.QueryRow(query, name)
	if err := row.Err(); err != nil {
		return zd.App{}, err
	}

	d, err := createFunctionFromRow(row)
	if err != nil {
		return zd.App{}, err
	}

	return d, nil
}

func (r *AppRepository) Create(function zd.App) (zd.App, error) {
	imagesStr := strings.Join(function.Images, ";")
	query := `INSERT INTO deployments(displayName, location, runtime, images) VALUES(?,?,?,?);`
	res, err := r.db.Exec(query, function.AppName, function.Location, function.Runtime, imagesStr)

	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return zd.App{}, zd.ErrRepoDuplicate
			}
		}

		return zd.App{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return zd.App{}, err
	}

	function.ID = id

	return function, nil
}

func (r *AppRepository) GetAll() ([]zd.App, error) {
	query := `SELECT * FROM deployments`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	deployments, err := createFunctionFromRows(rows)
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

// Helpers
// TODO: refactor this into one somehow
func createFunctionFromRows(rows *sql.Rows) ([]zd.App, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var deployments []zd.App

	for rows.Next() {
		var id int64
		var displayName string
		var location string
		var runtime string
		var images string

		err := rows.Scan(&id, &displayName, &location, &runtime, &images)
		if err != nil {
			return nil, err
		}

		d := zd.App{
			ID:       id,
			AppName:  displayName,
			Location: location,
			Runtime:  runtime,
			Images:   strings.Split(images, ";"),
		}

		deployments = append(deployments, d)
	}

	return deployments, nil
}

func createFunctionFromRow(row *sql.Row) (zd.App, error) {
	if err := row.Err(); err != nil {
		return zd.App{}, err
	}

	var id int64
	var displayName string
	var location string
	var runtime string
	var images string

	err := row.Scan(&id, &displayName, &location, &runtime, &images)
	if err != nil {
		return zd.App{}, err
	}

	d := zd.App{
		ID:       id,
		AppName:  displayName,
		Location: location,
		Runtime:  runtime,
		Images:   strings.Split(images, ";"),
	}

	return d, nil
}
