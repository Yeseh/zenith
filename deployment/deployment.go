package deployment

type Deployment struct {
	ID          int64
	DisplayName string
	Location    string
}

type CreateDeploymentDto struct {
	DisplayName string
	Location    string
}
