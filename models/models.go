package models

import "time"

type Workspace struct {
	Name     string
	Projects map[string]*Project
}

type Project struct {
	Name   string
	Timers []*Timer
}

type Timer struct {
	Labels []Label
	Start  time.Time
	End    time.Time
}

type Label struct {
	Name string
}
