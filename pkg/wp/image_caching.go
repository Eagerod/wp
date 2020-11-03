package wp

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

type ImageSource struct {
	SourcePath string
	LocalPath  string

	deleteParentDir bool
}

// Get the local image path from a given source image.
// This path will be a relative path that can be put anywhere needed.
func (is *ImageSource) GetImagePath() (string, error) {
	pathUrl, err := url.Parse(is.SourcePath)
	if err != nil {
		return "", err
	}

	if pathUrl.Scheme != "" && pathUrl.Scheme != "file" {
		return path.Join(pathUrl.Hostname(), filepath.Dir(pathUrl.Path)), nil
	}

	return "", nil
}

type FileDownloader func(destFile, sourceUrl string) error

// Ripped from https://golangcode.com/download-a-file-from-a-url
var downloadFile FileDownloader = func(destFile, sourceUrl string) error {
	resp, err := http.Get(sourceUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// Take the provided source path, and make a temporary copy of it that can be
//   fed through imagemagick repeatedly.
// CleaupImageSource must be called for the returned ImageSource.
func PrepareImageFromSource(sourcePath string, cacheDir string) (*ImageSource, error) {
	is := ImageSource{}
	is.SourcePath = sourcePath

	var outputDir string
	if cacheDir == "" {
		tempDir, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, err
		}

		outputDir = tempDir
		is.deleteParentDir = true
	} else {
		outputDir = cacheDir
		is.deleteParentDir = false
	}

	extraPath, err := is.GetImagePath()
	if err != nil {
		return nil, err
	}

	is.LocalPath = path.Join(outputDir, extraPath, path.Base(is.SourcePath))

	err = os.MkdirAll(filepath.Dir(is.LocalPath), 0755)
	if err != nil {
		return nil, err
	}

	if _, err = os.Stat(is.LocalPath); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err == nil {
		return &is, nil
	}

	pathUrl, err := url.Parse(is.SourcePath)
	if err != nil {
		return nil, err
	}

	if pathUrl.Scheme == "file" || pathUrl.Scheme == "" {
		source, err := os.Open(pathUrl.Path)
		if err != nil {
			return nil, err
		}
		defer source.Close()

		destination, err := os.Create(is.LocalPath)
		if err != nil {
			return nil, err
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)
	} else {
		err = downloadFile(is.LocalPath, is.SourcePath)
	}

	return &is, err
}

func CleanupImageSource(is *ImageSource) error {
	if is.deleteParentDir {
		return os.RemoveAll(path.Dir(is.LocalPath))
	}

	return nil
}
