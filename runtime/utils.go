package runtime

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

type TarBall struct {
	file *os.File
}

func TarFolder(source string) *bytes.Reader {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	dockerFile := "myDockerfile"
	dockerFileReader, err := os.Open("/path/to/dockerfile")
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	return dockerFileTarReader
}
