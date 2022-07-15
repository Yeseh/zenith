package runtime

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	zs "github.com/yeseh/zenith/mgmt/storage"
)

/* from create deployment route, the following should happen
1. Image is created for deployment
2. Image ID is returned
3. Deployment is stored in database with image ID
*/

type ImageConfig struct {
	Runtime            string
	AppName            string
	Location           string
	DockerBuildOptions types.ImageBuildOptions
}

func ImageNameForApp(appName string, runtime string) string {
	return fmt.Sprintf("zenith-apps/%s-%s", appName, runtime)
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

// TODO: TEMPORARY!!! Copies both the runtime definition, and the user provided app code into a temporary folder
// Then creates a TAR archive for the docker build context
func createBuildContext(runtime string, appName string) (io.ReadCloser, error) {
	wd, _ := os.Getwd()
	defPath := wd + "\\runtimes\\" + runtime

	ctxDir, err := ioutil.TempDir("", ".zenith-build-context")
	if err != nil {
		return nil, err
	}

	if err := zs.CopyDir(defPath, ctxDir); err != nil {
		return nil, err
	}

	appPath := zs.GetAppPath(appName)
	if err := zs.CopyDir(appPath, ctxDir); err != nil {
		return nil, err
	}

	ctx, err := archive.TarWithOptions(ctxDir, &archive.TarOptions{})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func CreateImageFor(config ImageConfig, docker DockerContext) (CreateImageResponse, error) {
	ctx, err := createBuildContext(config.Runtime, config.AppName)

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
