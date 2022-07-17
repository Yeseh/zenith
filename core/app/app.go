package app

import (
	"log"

	zd "github.com/yeseh/zenith/domain"
	zr "github.com/yeseh/zenith/runtime"
	zs "github.com/yeseh/zenith/storage"
)

type CreateAppDto struct {
	AppName        string `json:"appName"`
	Runtime        string `json:"runTime"`
	SourceLocation string `json:"sourceLocation"`
}

func ListApps(repo *AppRepository) ([]zd.App, error) {
	apps, err := repo.GetAll()
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func CreateApp(docker zr.DockerContext, repo *AppRepository, dto *CreateAppDto, storage zs.AppStorage) (zd.App, error) {
	target, err := storage.Upload(dto.SourceLocation, dto.AppName)
	if err != nil {
		log.Fatal(err)
	}

	config := zr.CreateImageConfig(dto.Runtime, dto.AppName, target)
	imageResponse, err := zr.CreateImage(config, docker, storage)
	if err != nil || !imageResponse.Success {
		log.Fatal(err)
	}

	function := zd.App{
		AppName:  dto.AppName,
		Runtime:  dto.Runtime,
		Location: target,
		Images:   config.DockerBuildOptions.Tags,
	}

	result, err := repo.Create(function)

	if err != nil {
		return zd.App{}, err
	}

	return result, nil
}
