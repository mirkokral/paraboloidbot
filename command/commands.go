package command

import "strings"

type Log struct {
	Info func(message string)
	Warn func(message string)
	Err  func(message string)
}

type Context struct {
	L        Log
	Args     []string
	Executor string
	Core     func(command string)
}

type Command struct {
	Name        string
	Description string
	Execute     func(c Context) string
}

var Commands = []Command{
	{
		Name:        "test",
		Description: "command and core test",
		Execute: func(c Context) string {
			c.L.Info("Hello world!")
			c.L.Info("&8Executed by: &a" + c.Executor)
			c.L.Info("&8Args: &7[&a\"" + strings.Join(c.Args, "&r&a\"&7, &a\"") + "&r&a\"&7]")
			return ""
		},
	},
}

func getCommands() []Command {
	return Commands
}
