package main

import (
	"fmt"
	"os"

	goarglib "github.com/vndg-rdmt/go-arglib"
)

func main() {

	// Command value holder, you man don't need this
	var command string

	// Place parser will put arg-values into
	var args struct {
		count int
		name  string
		not   bool
		arr   []string
		help  bool
	}

	parser := goarglib.New(
		goarglib.IntArg(&args.count, "-c", goarglib.Desc{
			Name:        "Counter",
			Description: "Count till this value",
		}),
		goarglib.StringArg(&args.name, "-n", goarglib.Desc{
			Name:        "Name",
			Description: "Ouput filename",
			Required:    true,
		}),
		goarglib.FlagArg(&args.not, "-r", goarglib.Desc{
			Name:        "Reqursive",
			Description: "Will programm execute recursively",
			Required:    true,
		}),
		goarglib.SliceArg(&args.arr, "-i", goarglib.Desc{
			Name:        "Ip",
			Description: "Ip addr-s of remote machine to send logs to",
			Required:    true,
		}),
		goarglib.FlagArg(&args.help, "-h", goarglib.Desc{
			Description: "Help page",
		}),
	)

	parser.DefineManual("This utility helps to do something VERY special.")

	parser.DefineCommands(&command,
		goarglib.NewCommand("create", "creates doc"),
		goarglib.NewCommand("connect", "connects to remote desktop"),
	)

	if err := parser.ParseArgs(os.Args[1:]); args.help {
		parser.WriteHelp(os.Stdout)
	} else if err != nil {
		fmt.Print(err)
		parser.WriteHelp(os.Stdout)
	}
}
