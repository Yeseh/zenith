package runtime

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/docker/docker/api/types"
	zs "github.com/yeseh/zenith/storage"
)

func ImageNameForApp(appName string, runtime string) string {
	return fmt.Sprintf("zenith-apps/%s-%s", appName, runtime)
}

func ListImages(docker DockerContext) ([]types.ImageSummary, error) {
	var zenImages []types.ImageSummary
	images, err := docker.Client.ImageList(docker.Ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, img := range images {
		for _, tag := range img.RepoTags {
			isZen := strings.HasPrefix(tag, "zenith-apps")
			if isZen {
				zenImages = append(zenImages, img)
			}
		}
	}

	return zenImages, nil
}

func CreateImageConfig(runtime string, appName string, path string) ImageConfig {
	imageName := ImageNameForApp(appName, runtime)
	tags := []string{
		// TODO: Properly version
		fmt.Sprintf("%s:%s", imageName, "latest"),
		// fmt.Sprintf("%s:%s", imageName, version),
	}

	opts := types.ImageBuildOptions{
		Remove: false,
		Tags:   tags,
	}

	return ImageConfig{
		Runtime:            runtime,
		AppName:            appName,
		Location:           path,
		DockerBuildOptions: opts,
	}
}

func CreateImage(config ImageConfig, docker DockerContext, storage zs.AppStorage) (CreateImageResponse, error) {
	ctx, err := storage.CreateContext(config.Runtime, config.AppName)

	if err != nil {
		return CreateImageResponse{}, err
	}

	buildResponse, err := docker.Client.ImageBuild(docker.Ctx, ctx, config.DockerBuildOptions)
	if err != nil {
		log.Fatal(err)
	}

	var outputBuf bytes.Buffer
	writer := io.Writer(&outputBuf)
	_, err = io.Copy(writer, buildResponse.Body)
	if err != nil {
		log.Fatal(err)
	}

	res := CreateImageResponse{
		Success:     true,
		BuildStream: string(outputBuf.Bytes()),
	}

	return res, nil
}
