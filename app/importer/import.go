package importer

import "io"

// Importer imports data.
type Importer interface {
	Read() ([]io.Reader, error)
}
