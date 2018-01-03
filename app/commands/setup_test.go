package commands

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/windler/ws/app/config"
	"github.com/windler/ws/internal/test"
)

func TestSetupCommand(t *testing.T) {
	f := new(SetupAppFactory).CreateCommand()

	assert.Equal(t, "setup", f.Command)
	assert.Equal(t, []string{}, f.Aliases)
	assert.Equal(t, "Configure everything to unleash the beauty. Alternatively, you can edit your personal config file.", f.Description)
	assert.Equal(t, 1, len(f.Subcommands))

	scWs := f.Subcommands[0]
	assert.Equal(t, "ws", scWs.Command)
	assert.Equal(t, []string{"workspace_dir"}, scWs.Aliases)
	assert.Equal(t, "Set the root dir where all of your workspaces are.", scWs.Description)
}

func TestSetNewWsDir(t *testing.T) {
	ui := test.MockUI()
	f := SetupAppFactory{
		UserInterface: ui,
	}

	c, _ := test.CreateTestContext(config.ConfigFlag)
	config.Repository().WsDir = ""

	oldStdin := mockStdIn("/testWsDir/")
	defer func() { os.Stdin = oldStdin }()

	err := f.CreateCommand().Subcommands[0].Action(c)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, "/testWsDir/", config.Repository().WsDir)
}

func mockStdIn(input string) *os.File {
	//https://stackoverflow.com/a/46365584
	content := []byte(input)
	tmpfile, err := ioutil.TempFile("", "wshero")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	tmpfile.Write(content)
	tmpfile.Seek(0, 0)

	oldStdin := os.Stdin
	os.Stdin = tmpfile

	return oldStdin
}
