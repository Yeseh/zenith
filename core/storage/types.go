package storage

import "io"

type AppStorage interface {
	Upload(source string, appName string) (string, error)
	Download(appName string) (string, error)
	CreateContext(runtime string, appName string) (io.ReadCloser, error)
	getAppPath(appName string) string
}
