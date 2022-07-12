package deployment

import (
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row doesn't exist")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type DeploymentRepository struct {
	db *sql.DB
}

func NewDeploymentRepository(db *sql.DB) *DeploymentRepository {
	return &DeploymentRepository{
		db,
	}
}

func (r *DeploymentRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS deployments(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		displayName TEXT NOT NULL UNIQUE,
		location TEXT NOT NULL
	);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *DeploymentRepository) Create(dto *CreateDeploymentDto) (*Deployment, error) {
	query := `INSERT INTO deployments(displayName, location) VALUES(?,?)`
	res, err := r.db.Exec(query, dto.DisplayName, dto.Location)

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

	deployment := &Deployment{
		ID:          id,
		DisplayName: dto.DisplayName,
		Location:    dto.Location,
	}

	return deployment, nil
}

func (r *DeploymentRepository) GetAll() ([]Deployment, error) {
	query := `SELECT id, displayName, location FROM deployments`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var deployments []Deployment

	for rows.Next() {
		var id int64
		var displayName string
		var location string

		err = rows.Scan(&id, &displayName, &location)
		if err != nil {
			return nil, err
		}

		d := Deployment{
			ID:          id,
			DisplayName: displayName,
			Location:    location,
		}

		deployments = append(deployments, d)
	}

	return deployments, nil
}
