package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
	"fmt"
	"github.com/crholm/timr/models"
	"time"
)

func ensure(err error){
	if err != nil {
		panic(err)
	}
}


func main() {

	loadConfig()
	loadStore()

	app := cli.NewApp()
	app.Name = "timr"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "english",
			Usage: "Language for the greeting",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}


	app.Commands = []cli.Command{
		{
			Name:    "stop",
			Usage:   "stops current activity",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "start",
			Usage:   "starts a new timer",
			ArgsUsage: "workspace project [label]",

			Action: start,
		},
		{
			Name:    "create",
			Usage:   "creates workspace and projects",
			Action: create,
		},



	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}




func create(c *cli.Context) error{
	if !(0 < len(c.Args())  && len(c.Args()) < 3 ){
		fmt.Println("There should only be 2 aguments <workspace> [project]")
		return nil
	}

	ws, err := store.LoadWorkspace(c.Args().Get(0))

	if err != nil {
		ws.Name = c.Args().Get(0)
	}

	if len(c.Args()) == 2 {
		name := c.Args()[1]

		if ws.Projects == nil{
			ws.Projects = make(map[string]*models.Project)
		}

		p := ws.Projects[name]
		p.Name = name
		ws.Projects[name] = p
	}

	store.SaveWorkspace(ws)

	return nil
}


func start(c *cli.Context) error {
	if len(c.Args()) != 2 {
		fmt.Println("There should only be 2 aguments <workspace> <project>")
		return nil
	}

	wsName := c.Args()[0]
	pName := c.Args()[1]


	ws, err := store.LoadWorkspace(wsName)
	if err != nil {
		fmt.Println("workspace does not exist")
		return nil
	}

	if ws.Projects == nil{
		fmt.Println("project does not exist")
	}

	p := ws.Projects[pName]

	if p.Name != pName {
		fmt.Println("project does not exist")
	}

	if p.Timers == nil{
		p.Timers = []*models.Timer{}
	}

	t := &models.Timer{
		Start: time.Now(),
	}

	p.Timers = append(p.Timers, t)

	err = store.SaveWorkspace(ws)

	if err != nil{
		fmt.Println("Could not save timer")
		return err
	}

	fmt.Println("Starting ", wsName, ">", pName, "@", t.Start.Format("2006-01-02 15:04:05 Z07:00"))


	return nil
}