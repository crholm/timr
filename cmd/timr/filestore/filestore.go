package filestore

import (
	"github.com/crholm/timr/models"
	"io/ioutil"
	"encoding/json"
)

type Filestore struct {
	Path string
}


func (f Filestore) LoadWorkspace(workspace string) (ws models.Workspace, err error){

	b, err := ioutil.ReadFile(f.Path + "/" + workspace)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &ws)
	return
}

func (f Filestore) SaveWorkspace(workspace models.Workspace) (err error){

	b, err := json.Marshal(workspace)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(f.Path + "/" + workspace.Name, b, 0644)
	return
}

func (f Filestore) ListWorkspace() []string{
	return []string{}
}