package init

type ZenithInitializer interface {
	Initialize(data InitData, output string) error
}

type InitData struct {
	AppName string
	Runtime string
}
