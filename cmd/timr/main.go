package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
	"fmt"
	"github.com/crholm/timr/models"
	"time"
	"github.com/golang-plus/errors"
	"strings"
)

func ensure(err error) {
	if err != nil {
		panic(err)
	}
}

const timeFormat = "2006-01-02 15:04 -07:00"

const runningFormat = "⏱ %s > %s @ %s"
const stopFormat = "⏹ %s > %s @ %s"


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
			Name:  "stop",
			Usage: "stops current activity",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "format, f",
					Value: stopFormat,
					Usage: "a format string for output",
				},
			},
			Action: stop,
		},
		{
			Name:      "start",
			Usage:     "starts a new timer",
			ArgsUsage: "workspace project [label]",
			Action: start,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "format, f",
					Value: runningFormat,
					Usage: "a format string for start output",
				},
				cli.StringFlag{
					Name: "stop-format, s",
					Value: stopFormat,
					Usage: "a format string for stop output",
				},
			},


		},
		{
			Name:   "create",
			Usage:  "creates workspace and projects",
			Action: create,
		},
		{
			Name:   "status",
			Usage:  "Lists running timers",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "format, f",
					Value: runningFormat,
					Usage: "a format string for output",
				},
			},
			Action: status,
		},
		{
			Name:   "list",
			Usage:  "Lists workspaces, projects or logged items",
			Action: list,
		},
		{
			Name:   "export",
			Usage:  "Export data as tsv",
			Action: export,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func formatDuration(d time.Duration) string{
	d = d.Truncate(time.Second)
	s := d.String()
	s = strings.Replace(s, "h", "h ", 1)
	s = strings.Replace(s, "m", "m ", 1)
	return s
}

func create(c *cli.Context) error {
	if !(0 < len(c.Args()) && len(c.Args()) < 3) {
		fmt.Println("There should only be 2 aguments <workspace> [project]")
		return nil
	}

	ws, err := store.LoadWorkspace(c.Args().Get(0))

	if err != nil {
		ws.Name = c.Args().Get(0)
	}

	if len(c.Args()) == 2 {
		name := c.Args()[1]

		if ws.Projects == nil {
			ws.Projects = make(map[string]*models.Project)
		}

		p := ws.Projects[name]

		if p == nil {
			p = &models.Project{}
		}

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

	printFormat := c.String("format")
	printStopFormat := c.String("stop-format")


	wsName := c.Args()[0]
	pName := c.Args()[1]

	startTime := time.Now()
	z := time.Time{}

	ws, err := store.LoadWorkspace(wsName)
	if err != nil {
		fmt.Println("workspace does not exist")
		return nil
	}

	for _, project := range ws.Projects{

		if len(project.Timers) == 0{
			continue
		}

		timer := project.Timers[len(project.Timers) - 1]
		if timer.End == z{
			timer.End = startTime
			fmt.Printf(printStopFormat,  ws.Name , project.Name, formatDuration(timer.End.Sub(timer.Start)))
			fmt.Println()
		}

	}


	if ws.Projects == nil {
		fmt.Println("project does not exist")
		return nil
	}


	p := ws.Projects[pName]

	if p == nil {
		fmt.Println("project does not exist")
		return nil
	}

	if p.Timers == nil {
		p.Timers = []*models.Timer{}
	}

	t := &models.Timer{
		Start: time.Now(),
	}

	p.Timers = append(p.Timers, t)

	err = store.SaveWorkspace(ws)

	if err != nil {
		fmt.Println("Could not save timer")
		return err
	}

	fmt.Printf(printFormat, wsName, pName, t.Start.Format(timeFormat))
	fmt.Println()

	return nil
}


func stop(c *cli.Context) error {

	printFormat := c.String("format")

	end := time.Now()
	z := time.Time{}
	doStop := func(name string){
		ws, err := store.LoadWorkspace(name)
		ensure(err)

		for _, project := range ws.Projects{

			if len(project.Timers) == 0{
				continue
			}

			timer := project.Timers[len(project.Timers) - 1]
			if timer.End == z{
				timer.End = end
				fmt.Printf(printFormat, name, project.Name, timer.End.Sub(timer.Start).Truncate(time.Second))
				fmt.Println()
			}

		}
		store.SaveWorkspace(ws)

	}

	if len(c.Args()) == 0 {
		workspaces, err := store.ListWorkspace()
		ensure(err)
		for _, workspace := range workspaces {
			doStop(workspace)
		}
	}

	if len(c.Args()) == 1 {
		doStop(c.Args().Get(0))
	}


	return nil
}


func status(c *cli.Context) error {

	printFormat := c.String("format")

	z := time.Time{}
	printStatus := func(name string){
		ws, err := store.LoadWorkspace(name)
		ensure(err)

		for _, project := range ws.Projects{

			if len(project.Timers) == 0{
				continue
			}

			timer := project.Timers[len(project.Timers) - 1]
			if timer.End == z{
				fmt.Printf(printFormat, name, project.Name, formatDuration(time.Now().Sub(timer.Start)) )
				fmt.Println()
			}

		}
	}

	if len(c.Args()) == 0 {
		workspaces, err := store.ListWorkspace()
		ensure(err)
		for _, workspace := range workspaces {
			printStatus(workspace)
		}
	}

	if len(c.Args()) == 1 {
		printStatus(c.Args().Get(0))
	}


	return nil
}


func list(c *cli.Context) error {

	if len(c.Args()) == 0 {
		workspaces, err := store.ListWorkspace()
		ensure(err)
		sort.Strings(workspaces)
		for _, w := range workspaces {
			fmt.Println(w)
		}
		return nil
	}

	ws, err := store.LoadWorkspace(c.Args().Get(0))
	ensure(err)

	if len(c.Args()) == 1 {
		var names []string
		for project, _ := range ws.Projects {
			names = append(names, project)
		}
		sort.Strings(names)
		for _, project := range names {
			fmt.Println(project)
		}
		return nil
	}

	if len(c.Args()) == 2 {
		project := ws.Projects[c.Args().Get(1)]
		if project == nil{
			ensure(errors.New("Could not find Project"))
		}

		z := time.Time{}
		for i, timer := range project.Timers{


			fmt.Print(i)
			fmt.Print("\t")
			fmt.Print(timer.Start.Format(timeFormat))

			if timer.End != z {
				fmt.Print("\t")
				fmt.Print(timer.End.Format(timeFormat))
				fmt.Print("\t")
				fmt.Print(formatDuration(timer.End.Sub(timer.Start)))
			}else {
				fmt.Print("\t\t")
			}

			fmt.Print("\n")


		}


	}

	return nil
}


func export(c *cli.Context) error {


	z := time.Time{}
	exportTsv := func(name string){
		ws, err := store.LoadWorkspace(name)
		ensure(err)

		for _, project := range ws.Projects{
			for _, timer := range project.Timers{
				fmt.Print(ws.Name)
				fmt.Print("\t")
				fmt.Print(project.Name)
				fmt.Print("\t")
				fmt.Print(timer.Start)

				if timer.End != z {
					fmt.Print("\t")
					fmt.Print(timer.End)
					fmt.Print("\t")
					fmt.Print(timer.End.Sub(timer.Start).Hours())
				}else {
					fmt.Print("\t\t")
				}

				fmt.Println()
			}
		}
	}

	if len(c.Args()) == 0 {
		workspaces, err := store.ListWorkspace()
		ensure(err)
		for _, workspace := range workspaces {
			exportTsv(workspace)
		}
	}

	if len(c.Args()) == 1 {
		exportTsv(c.Args().Get(0))
	}


	return nil
}





