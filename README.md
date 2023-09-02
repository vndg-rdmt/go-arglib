# GO-ARGLIB

Simple minimalistic and very obvious in usage.

## How to use


Create new parser, and define flags/argument.
Supported types are string, int, flag(bool) and array(slice). Map is in feature, as passing array via braces.

```go
var count int
var name string

parser := goarglib.New(
    goarglib.IntArg(&count, "-c", goarglib.Desc{
	Name:        "Counter",
	Description: "Count till this value",
    }),
    goarglib.StringArg(&name, "-n", goarglib.Desc{
	Name:        "Name",
	Description: "Ouput filename",
	Required:    true,
    }),
)
```

You can also define commands

```go
var command string

parser.DefineCommands(&command,
    goarglib.NewCommand("create", "creates doc"),
    goarglib.NewCommand("connect", "connects to remotedesktop"),
)
```


Parsing arguments. You can pass any slice of strings for parser, so you can use other than os.Args sources for reading arguments

```go
var help bool

parser := goarglib.New(
    goarglib.IntArg(&count, "-h", goarglib.Desc{
	Description: "Help page",
    }),
)

parser.DefineManual("This utility helps to do something VERY special.")

if err := parser.ParseArgs(os.Args[1:]); help {
    parser.WriteHelp(os.Stdout)

} else if err != nil {
    fmt.Print(err)
    parser.WriteHelp(os.Stdout)
}
```

Output example.
> Source code can be found in example diretory. Use make command to run.

```text
Usage: app COMMAND [OPTIONS]
This utility helps to do something VERY special.

Commands:
  create   creates doc
  connect  connects to remote desktop

Options:
  -i   array Ip          Ip addr-s of remote machine to send logs to
  -c  number Counter     Count till this value
  -n  string Name        Ouput filename
  -r    flag Reqursive   Will programm execute recursively
```

#### Errors handling

`go-arglib` never panics, but instead - returns errors.

This situations will return errors:
- Passing wrong type as an argument, for example passed array to string type, or for flag types (bool) passed any external symbols
- Not all params, marked as `Required: true`, are passed
- Catched undefined param
- The same argument key passed twice (or more)


#### Man page generation

Man page sections are generated dynamicaly, so you want seen something like Usage: app COMMANDS, while you don't defined it. App name is defined from executable name, so you don't need to change it, it will always have the actual name. Name and descriptos for keys are not neccesary, just the key and a pointer.

Result can be seen by running example/example.go