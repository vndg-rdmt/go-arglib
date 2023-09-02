package goarglib

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vndg-rdmt/go-arglib/internal/fsm"
)

type parser struct {
	manual      string
	commands    map[string]command
	keysMapping map[string]argument
}

// Constructor for `parser`.
// Receives `arguments` as an input, which can be created
// with package Arg functions for a specific type
func New(args ...argument) *parser {
	buf := &parser{
		keysMapping: make(map[string]argument, len(args)),
	}
	for _, arg := range args {
		buf.keysMapping[arg.key] = arg
	}
	return buf
}

// Defines short-manual for programm
func (self *parser) DefineManual(text string) {
	self.manual = text
}

// Adds command definitions to arglib parser.
// Will overwrite previously defined commands with
// new ones.
func (self *parser) DefineCommands(p *string, comm ...command) {
	setter := func(value string) { *p = value }
	self.commands = make(map[string]command)

	for _, v := range comm {
		v.setValue = setter
		self.commands[v.key] = v
	}
}

// Parses arguments from the provided source.
//
// Typically, argument parsers provide an ability to lookup
// only os.Args, but arglib goes above and provides you
// an ability to parse any source as others do to os.Args.
//
// If you want to achive "default" behavior, pass os.Args[1:]
// as an argument. (range because of the fact that 0 arg - executable name).
func (self *parser) ParseArgs(source []string) error {
	if len(source) == 0 {
		return nil
	}

	var searchIndex int

	if self.commands != nil && len(self.commands) > 0 {
		v := source[searchIndex]
		if comm, ok := self.commands[v]; ok {
			comm.setValue(v)
			searchIndex++
		}
	}

	fsmInst := fsm.New(&source)

	for ; searchIndex < len(source); searchIndex++ {
		if _, ok := self.keysMapping[source[searchIndex]]; ok {
			if err := fsmInst.CloseState(searchIndex); err != nil {
				return err
			}
			fsmInst.SwitchState(searchIndex)
		} else {
			if !fsmInst.IsValidToProcced() {
				return fmt.Errorf("Found unrecognized parameter: %s\n", source[searchIndex])
			}
		}
	}

	parsingResult := *fsmInst.Finally()
	missingParams := make([]string, 0, len(self.keysMapping))

	for key, arg := range self.keysMapping {
		if value, ok := parsingResult[key]; ok {
			if err := arg.setValue(value); err != nil {
				return err
			}
		} else if arg.desc.Required {
			missingParams = append(missingParams, key)
		}
	}

	if len(missingParams) > 0 {
		return fmt.Errorf("Missing required params: %s\n", strings.Join(missingParams, " "))
	}
	return nil
}

// Generates help/manual page based on defined
// commands and arguments.
func (self *parser) GenerateHelp() string {
	headerBuffer := strings.Builder{}
	buffer := strings.Builder{}

	headerBuffer.WriteString("\nUsage: ")
	headerBuffer.WriteString(filepath.Base(os.Args[0]))

	if self.commands != nil && len(self.commands) > 0 {
		headerBuffer.WriteString(" COMMAND")
		buffer.WriteString("\nCommands:\n")

		var nameLength int = 0
		for _, v := range self.commands {
			if len(v.key) > nameLength {
				nameLength = len(v.key)
			}
		}

		for _, v := range self.commands {
			buffer.WriteString(
				fmt.Sprintf("  %s%s  %s\n",
					v.key,
					strings.Repeat(" ", nameLength-len(v.key)),
					v.description,
				),
			)
		}
	}

	if len(self.keysMapping) > 0 {
		headerBuffer.WriteString(" [OPTIONS]")
		buffer.WriteString("\nOptions:\n")

		var keyLength int = 0
		var nameLength int = 0
		for _, v := range self.keysMapping {
			if len(v.key)+len(v.valueType) > keyLength {
				keyLength = len(v.key) + len(v.valueType)
			}
			if len(v.desc.Name) > nameLength {
				nameLength = len(v.desc.Name)
			}
		}

		for _, v := range self.keysMapping {
			buffer.WriteString(
				fmt.Sprintf("  %s %s %s %s%s   %s\n",
					v.key,
					strings.Repeat(" ", keyLength-len(v.key)-len(v.valueType)),
					v.valueType,
					v.desc.Name,
					strings.Repeat(" ", nameLength-len(v.desc.Name)),
					v.desc.Description,
				),
			)
		}
	}

	if self.manual != "" {
		headerBuffer.WriteString("\n")
		headerBuffer.WriteString(self.manual)
	}
	headerBuffer.WriteString("\n")
	headerBuffer.WriteString(buffer.String())
	headerBuffer.WriteString("\n")

	return headerBuffer.String()
}

// Writes generated help/manual page to the provided io.Writer.
// Returns writer.Write method call type - usually, the number
// of bytes written and an error, if any.
func (self *parser) WriteHelp(writer io.Writer) (int, error) {
	return writer.Write([]byte(self.GenerateHelp()))
}
