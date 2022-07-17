package domain

import "errors"

var (
	// Repository errors
	ErrRepoDuplicate    = errors.New("record already exists")
	ErrRepoNotExists    = errors.New("row doesn't exist")
	ErrRepoUpdateFailed = errors.New("update failed")
	ErrRepoDeleteFailed = errors.New("delete failed")
)
