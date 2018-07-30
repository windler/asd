package commands

import (
	"bytes"
	"html/template"
	"os"
	"sort"
	"strings"
)

//ListWsFactory creates commands to list workspace information
type ListWsFactory struct {
	InfoRetriever WsInfoRetriever
	UserInterface UI
	WSRetriever   WorkspaceRetriever
	Executor      CustomCommandExecutor
}

type WsInfoRetriever interface {
	Status(ws string) string
	CurrentBranch(ws string) string
}

type tableData [][]string

//CreateCommand creates a ListWsCommand
func (factory *ListWsFactory) CreateCommand() BaseCommand {

	return BaseCommand{
		Command:     "ls",
		Description: "List all workspaces with fancy information.",
		Aliases:     []string{},
		Action: func(c WSCommandContext) {
			factory.listWsExec(&c)
		},
		Flags: []StringFlag{
			StringFlag{
				"table",
				"formats the table using the `template`",
			},
		},
	}
}

func (factory *ListWsFactory) UI() UI {
	return factory.UserInterface
}

func (c tableData) Len() int           { return len(c) }
func (c tableData) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c tableData) Less(i, j int) bool { return strings.Compare(c[i][0], c[j][0]) == -1 }

func (factory *ListWsFactory) listWsExec(c *WSCommandContext) {
	conf := (*c).GetConfig()

	wsDir := conf.GetWsDir()

	if wsDir == "" {
		factory.UI().PrintString("Panic!", "red")
		factory.UI().PrintString(" >> No workspaces defined to scan <<")

		return
	}

	dirs := factory.WSRetriever.GetWorkspacesIn(wsDir)

	dataChannel := factory.channelFileInfos(dirs)
	fanOutChannels := []<-chan []string{}

	tableFormat := getTableFormat(c)

	parallel := 3
	if conf.GetParallelProcessing() > 0 {
		parallel = conf.GetParallelProcessing()
	}

	for i := 0; i < parallel; i++ {
		fanOutChannels = append(fanOutChannels, factory.collectWsData(dataChannel, tableFormat, c))
	}

	rows := tableData{}
	fanInChannel := factory.fanIn(fanOutChannels)
	for i := 0; i < len(dirs); i++ {
		rows = append(rows, <-fanInChannel)
	}

	if len(rows) > 0 {
		sort.Sort(rows)

		funcMap := factory.getHeaderFunctionMap()

		buf := new(bytes.Buffer)
		t := template.Must(template.New("header").Funcs(funcMap).Parse(tableFormat))
		t.Execute(buf, "")

		factory.UI().PrintTable(strings.Split(buf.String(), "|"), rows)
	} else {
		factory.printError()
	}
}

func getTableFormat(c *WSCommandContext) string {
	conf := (*c).GetConfig()

	tableFormat := "{{wsRoot .}}|{{gitStatus .}}|{{gitBranch .}}"
	if (*c).GetStringFlag("table") != "" {
		tableFormat = (*c).GetStringFlag("table")
	} else if conf.GetTableFormat() != "" {
		tableFormat = conf.GetTableFormat()
	}

	return tableFormat
}

func (factory *ListWsFactory) printError() {
	factory.UI().PrintString("No workspaces found!", "red")
}

func (factory *ListWsFactory) channelFileInfos(dirs []string) <-chan string {
	out := make(chan string, len(dirs))
	go func() {
		for _, dir := range dirs {
			out <- dir
		}
		close(out)
	}()
	return out
}

func (factory *ListWsFactory) collectWsData(in <-chan string, pattern string, c *WSCommandContext) <-chan []string {
	out := make(chan []string)
	go func() {
		funcMap := factory.getRowsFunctionMap(true, c)

		for dir := range in {
			buf := new(bytes.Buffer)
			t := template.Must(template.New("table").Funcs(funcMap).Parse(pattern))
			t.Execute(buf, dir)

			out <- strings.Split(buf.String(), "|")
		}
		close(out)
	}()
	return out
}

type tableTemplateData struct {
	Dir string
}

func (factory *ListWsFactory) fanIn(input []<-chan []string) <-chan []string {
	c := make(chan []string)
	for _, ch := range input {
		go func(channel <-chan []string) {
			for msg := range channel {
				c <- msg
			}
		}(ch)
	}

	return c
}

func (factory *ListWsFactory) getHeaderFunctionMap() template.FuncMap {
	return template.FuncMap{
		"wsRoot":    func(dir string) string { return "ws" },
		"gitStatus": func(dir string) string { return "git status" },
		"gitBranch": func(dir string) string { return "git branch" },
		"cmd":       func(name, dir string) string { return name },
	}
}

func (factory *ListWsFactory) getRowsFunctionMap(markCurrentWs bool, c *WSCommandContext) template.FuncMap {
	return template.FuncMap{
		"wsRoot": func(dir string) string {
			res := dir
			wd, _ := os.Getwd()
			if markCurrentWs && strings.HasPrefix(wd, dir) {
				res = res + " <--"
			}
			return res
		},
		"gitStatus": func(dir string) string {
			return factory.InfoRetriever.Status(dir)
		},
		"gitBranch": func(dir string) string {
			return factory.InfoRetriever.CurrentBranch(dir)
		},
		"cmd": func(name, dir string) string {
			for _, cmd := range (*c).GetConfig().GetCustomCommands() {
				if cmd.GetName() == name {
					return strings.TrimSpace(factory.Executor.ExecToString(&cmd, dir, c))
				}
			}
			return "-- NO OUTPUT --"
		},
	}
}
