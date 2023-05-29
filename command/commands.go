package command

import "strings"

type Log struct {
	Info func(message string)
	Warn func(message string)
	Err  func(message string)
}

type Command struct {
	Name        string
	Description string
	Execute     func(l Log, args []string, executor string) string
}

var Commands = []Command{
	{
		Name:        "test",
		Description: "command and core test",
		Execute: func(l Log, args []string, executor string) string {
			l.Info("Hello world!")
			l.Info("&8Executed by: &a" + executor)
			l.Info("&8Args: &7[&a\"" + strings.Join(args, "&r&a\"&7, &a\"") + "&r&a\"&7]")
			return ""
		},
	},
}

func getCommands() []Command {
	return Commands
}
