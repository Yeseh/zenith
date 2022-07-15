package storage

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

const ZENITH_MOUNT string = "D:/ZenithMount"

// TODO: VERY TEMPORARY!! Create interface 'AppUploader' for uploading to Azure storage etc
// For now just copy local files to local zenithmount folder

func UploadApp(source string, appName string) (string, error) {
	target := path.Join(ZENITH_MOUNT, "/apps", appName)

	cmd := exec.Command("powershell", "Remove-Item", "-Recurse", "-Force", target)
	if err := cmd.Run(); err != nil {
		return target, err
	}

	cmd = exec.Command("powershell", "Copy-Item", "-Recurse", source, target)
	err := cmd.Run()

	return target, err
}

func DownloadApp(appName string) (string, error) {
	source := GetAppPath(appName)
	target, err := ioutil.TempDir("", ".app")
	if err != nil {
		return target, err
	}

	if err := CopyDir(source, target); err != nil {
		return target, err
	}

	return target, nil
}

func GetAppPath(appName string) string {
	return path.Join(ZENITH_MOUNT, "/apps", appName)
}

// https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
func CopyFile(src, dst string) (err error) {
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

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
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
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
