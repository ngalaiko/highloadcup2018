package app

import (
	"github.com/ngalayko/highloadcup/app/datastore"
	"github.com/ngalayko/highloadcup/app/importer/zip"
	"github.com/ngalayko/highloadcup/app/logger"
	"github.com/ngalayko/highloadcup/app/web"
)

// Application is a main object.
type Application struct {
	datastore *datastore.Datastore
	web       *web.Web
}

// New is the application constructor.
func New(dataPath string) (*Application, error) {
	logger := logger.New()

	datastore, err := datastore.New(logger, zip.New(dataPath))
	if err != nil {
		return nil, err
	}

	return &Application{
		web:       web.New(logger, datastore),
		datastore: datastore,
	}, nil
}

// ListenAndServe starts the server.
func (a *Application) ListenAndServe(addr string) error {
	return a.web.ListenAndServe(addr)
}
