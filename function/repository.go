package function

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row doesn't exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type FunctionRepository struct {
	db *sql.DB
}

func NewFunctionRepository(db *sql.DB) *FunctionRepository {
	return &FunctionRepository{
		db,
	}
}

func (r *FunctionRepository) Migrate() error {
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

func createFunctionFromRows(rows *sql.Rows) ([]Function, error) {
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var deployments []Function

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

		d := Function{
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

func createFunctionFromRow(row *sql.Row) (Function, error) {
	if err := row.Err(); err != nil {
		return Function{}, err
	}

	var id int64
	var displayName string
	var location string
	var runtime string
	var images string

	err := row.Scan(&id, &displayName, &location, &runtime, &images)
	if err != nil {
		return Function{}, err
	}

	d := Function{
		ID:       id,
		AppName:  displayName,
		Location: location,
		Runtime:  runtime,
		Images:   strings.Split(images, ";"),
	}

	return d, nil
}

func (r *FunctionRepository) GetByDisplayName(name string) (Function, error) {
	query := `SELECT * FROM deployments WHERE displayName = ?`
	row := r.db.QueryRow(query, name)
	if err := row.Err(); err != nil {
		return Function{}, err
	}

	d, err := createFunctionFromRow(row)
	if err != nil {
		return Function{}, err
	}

	return d, nil
}

func (r *FunctionRepository) Create(function Function) (*Function, error) {
	imagesStr := strings.Join(function.Images, ";")
	query := `INSERT INTO deployments(displayName, location, runtime, images) VALUES(?,?,?,?);`
	res, err := r.db.Exec(query, function.AppName, function.Location, function.Runtime, imagesStr)

	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}

		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	function.ID = id

	return &function, nil
}

func (r *FunctionRepository) GetAll() ([]Function, error) {
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
