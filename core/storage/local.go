package storage

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
)

const ZENITH_MOUNT = "D:/ZenithMount"

// Storage manager for local development. Makes use of a 'normal' filesystem to store Zenith applications;
type LocalStorage struct {
	mountDir string
}

func (l LocalStorage) New(zenithMount string) *LocalStorage {
	return &LocalStorage{
		mountDir: zenithMount,
	}
}

func (l LocalStorage) NewFromEnv(zenithMount string) (*LocalStorage, error) {
	mount := os.Getenv("ZENITH_MOUNT")
	if len(mount) == 0 {
		return nil, errors.New("ZENITH_MOUNT environment variable not set")
	}

	cli := LocalStorage{
		mountDir: mount,
	}

	return &cli, nil
}

// Uploads an application to the local mount directory;
func (l *LocalStorage) Upload(source string, appName string) (string, error) {
	target := l.getAppPath(appName)

	if err := os.RemoveAll(target); err != nil {
		return target, err
	}

	if err := copyDir(source, target); err != nil {
		return target, err
	}

	return target, nil
}

// Copy app with appname to a temporary directory.
// Caller is responsible for cleaning up the temporary dir.
func (l *LocalStorage) DownloadApp(appName string) (string, error) {
	source := l.getAppPath(appName)
	target, err := ioutil.TempDir("", ".app")
	if err != nil {
		return target, err
	}

	if err := copyDir(source, target); err != nil {
		return target, err
	}

	return target, nil
}

// Creates a TAR archive from the user provided code, and the Zenith runtime
// This can be used as the buildcontext for Docker images
func (l *LocalStorage) CreateContext(runtime string, appName string) (io.ReadCloser, error) {
	wd, _ := os.Getwd()
	rtPath := wd + "\\_runtimes\\" + runtime

	ctxDir, err := ioutil.TempDir("", ".zenith-build-context")
	if err != nil {
		return nil, err
	}

	// Clean up temp dir after getting TAR'd
	defer os.RemoveAll(ctxDir)
	if err := copyDir(rtPath, ctxDir); err != nil {
		return nil, err
	}

	appPath := l.getAppPath(appName)
	if err := copyDir(appPath, ctxDir); err != nil {
		return nil, err
	}

	ctx, err := archive.TarWithOptions(ctxDir, &archive.TarOptions{})
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

func (l *LocalStorage) getAppPath(appName string) string {
	return path.Join(l.mountDir, "/apps", appName)
}

// https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	// os.Stat returns error when not exists
	// Create destination folder if not exists
	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
