package goarglib

type command struct {
	key         string
	description string
	setValue    func(value string)
}

// Creates new command. Does not require pointer due
// to it's defined once within parser.DefineCommands function call.
func NewCommand(key string, description string) command {
	return command{
		key:         key,
		description: description,
	}
}
