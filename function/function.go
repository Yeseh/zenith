package function

type Function struct {
	ID       int64
	AppName  string
	Location string
	Runtime  string
	Images   []string
}

type CreateFunctionDto struct {
	AppName        string `json:"appName"`
	Runtime        string `json:"runTime"`
	SourceLocation string `json:"sourceLocation"`
}
