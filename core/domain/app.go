package domain

type App struct {
	ID       int64
	AppName  string
	Location string
	Runtime  string
	Images   []string
}

type AppRepository interface {
	GetAll() ([]App, error)
	GetByDisplayName(name string) (App, error)
	Create() (App, error)
}
