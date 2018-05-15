package main

import (
	"github.com/crholm/timr/cmd/timr/filestore"
	"github.com/crholm/timr/models"
	"os"
)

type Store interface {
	LoadWorkspace(workspace string) (models.Workspace, error)
	SaveWorkspace(workspace models.Workspace) (error)
    ListWorkspace() ([]string, error)
}


var store Store

func loadStore(){

	switch config.Store {
	case "file":
		path := configDir+"/filestore"
		store = filestore.Filestore{Path: path}
		err := os.MkdirAll(path,os.ModePerm)
		ensure(err)

	}
}