package wp

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ImageSource struct {
	SourcePath string
	LocalPath  string

	deleteParentDir bool
}

func (is *ImageSource) SourcePathIsLocal() (bool, error) {
	pathUrl, err := url.Parse(is.SourcePath)
	if err != nil {
		return false, err
	}

	return pathUrl.Scheme == "" || pathUrl.Scheme == "file", nil
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

	// If the SourcePath is rooted in the cache directory, bail early, because
	//   there's nothing interesting to do, and trying to use any kind of
	//   recaching logic will just duplicate the image in the cache directory.
	// Doesn't obey obey symlinks, for now.
	if cacheDir != "" {
		isLocal, err := is.SourcePathIsLocal()
		if err != nil {
			return nil, err
		}

		if isLocal {
			absSource, err := filepath.Abs(is.SourcePath)
			if err != nil {
				return nil, err
			}

			absCache, err := filepath.Abs(cacheDir)
			if err != nil {
				return nil, err
			}

			if strings.HasPrefix(absSource, absCache) {
				if _, err := os.Stat(is.SourcePath); err != nil {
					return nil, err
				}

				is.LocalPath = is.SourcePath
				is.deleteParentDir = false
				return &is, nil
			}
		}
	}

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
