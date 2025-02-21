package models

type CommandCallback func(login string, args []string)

type ChatCommand struct {
	Name     string
	Callback CommandCallback
	Admin    bool
	Help     string
}