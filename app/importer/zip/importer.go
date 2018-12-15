package zip

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/ngalayko/highloadcup/app/logger"
)

// Importer imports data from zip file.
type Importer struct {
	logger logger.Logger
	path   string
}

// New creates zip importer.
func New(path string) *Importer {
	return &Importer{
		path: path,
	}
}

// Read returns data readers from zip archive.
func (z *Importer) Read() ([]io.Reader, error) {
	rc, err := zip.OpenReader(z.path)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	result := make([]io.Reader, 0, len(rc.File))
	for _, f := range rc.File {
		if filepath.Ext(f.Name) != ".json" {
			z.logger.Info("skipping file %s", f.Name)
			continue
		}

		z.logger.Info("importing file %s", f.Name)

		frc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("can't open file '%s': %s", f.Name, err)
		}

		data, err := ioutil.ReadAll(frc)
		if err != nil {
			return nil, fmt.Errorf("can't read file '%s': %s", f.Name, err)
		}

		result = append(result, bytes.NewBuffer(data))

		frc.Close()
	}
	return result, nil
}
