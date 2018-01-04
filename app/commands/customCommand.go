package commands

import (
	"fmt"

	"github.com/windler/ws/app/commands/internal/commandCommons"

	"github.com/fatih/color"

	"github.com/urfave/cli"
	"github.com/windler/ws/app/config"
)

//SetupAppFactory creates commands to list workspace information
type CustomCommandFactory struct {
	UserInterface UI
	Cmd           config.CustomCommand
}

//ensure interface
var _ BaseCommandFactory = &CustomCommandFactory{}

//CreateCommand creates a ListWsCommand
func (factory *CustomCommandFactory) CreateCommand() BaseCommand {
	return BaseCommand{
		Command:     factory.Cmd.Name,
		Description: factory.Cmd.Description,
		Action: func(c *cli.Context) error {
			return factory.action(c)
		},
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "w, ws",
				Usage: "executes the command in `workspace` instead of the current workspace. `workspace` is a regex pattern that must match a single workspace.",
			},
		},
	}
}

func (factory *CustomCommandFactory) UI() UI {
	return factory.UserInterface
}

func (factory *CustomCommandFactory) action(c *cli.Context) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Custom command is invalid. Check config.", r)
		}
	}()

	factory.UI().PrintString(factory.Cmd.Name+":", color.FgGreen)

	var output string
	if c.String("ws") != "" {
		ws := commandCommons.GetWorkspaceByPattern(config.Repository().WsDir, c.String("ws"))
		output = commandCommons.ExecCustomCommand(&factory.Cmd, ws)
	} else {
		output = commandCommons.ExecCustomCommandInCurrentWs(&factory.Cmd)
	}

	factory.UI().PrintString(output)

	return nil
}
