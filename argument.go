package goarglib

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	typeString = "string"
	typeNumber = "number"
	typeFlag   = "flag"
	typeArray  = "array"
)

// Argument descriptor.
// Struct is used to describe argument flag and its
// behavior regardless of value type.
type Desc struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
}

// Argument descriptor decorator, which uses `setValue`
// as a type-regardless argument value type converter
// and setter.
//
// Struct also holds argument descriptor to provide
// unified interface for `parser` to go through arguments
// string, parse it and describe simultaneously.
type argument struct {
	desc      Desc
	key       string
	valueType string
	setValue  func(value []string) error
}

// @unsafe-architecture
// type argValueHandlerConstructor[T any] func(p *T) argValueHandler
// type argumentConstructor[T any] func(p *T, key string, descriptor Desc) argument
// func newArgumentConstructor[ArgType any](cons argValueHandlerConstructor[ArgType]) argumentConstructor[ArgType] {
// 	return func(pt *ArgType, key string, desc Desc) argument {
// 		return argument{
// 			desc, key, cons(pt),
// 		}
// 	}
// }

// Creates a new `argument` for parser for a string arg type.
func StringArg(p *string, key string, descriptor Desc) argument {
	return argument{
		descriptor,
		key,
		typeString,
		func(value []string) error {
			if len(value) == 1 {
				*p = value[0]
				return nil
			}
			return genArgTypeError(key, typeString, &value)
		},
	}
}

// Creates a new `argument` for parser for an integer arg type.
func IntArg(p *int, key string, descriptor Desc) argument {
	return argument{
		descriptor,
		key,
		"number",
		func(value []string) error {
			if len(value) == 1 {
				if res, err := strconv.Atoi(value[0]); err == nil {
					*p = res
					return nil
				}
			}
			return genArgTypeError(key, typeNumber, &value)
		},
	}
}

// Creates a new `argument` for parser for a boolean arg type.
func FlagArg(p *bool, key string, descriptor Desc) argument {
	return argument{
		descriptor,
		key,
		typeFlag,
		func(value []string) error {
			if len(value) == 0 {
				*p = true
				return nil
			}
			return genArgTypeError(key, typeFlag, &value)
		},
	}
}

// Creates a new `argument` for parser for an array arg type.
// Currently supported only a []string type, but it will be replaced
// in future with unified functions.
func SliceArg(p *[]string, key string, descriptor Desc) argument {
	return argument{
		descriptor,
		key,
		"array",
		func(value []string) error {
			*p = value
			return nil
		},
	}

}

func genArgTypeError(key string, valueType string, passed *[]string) error {
	return fmt.Errorf("Argument flag '%s' requires value of type '%s', got [%s] \n", key, valueType, strings.Join(*passed, ", "))
}

// @featured
// Creates a new `argument` for parser for a map arg type.
// Currently supported only a map[string]string type, but it will be replaced
// in future with unified functions.
// func MapArg(p *map[string]string, key string, descriptor Desc) argument[map[string]string] {
// 	return argument[map[string]string]{
// 		descriptor,
// 		key,
// 		"map",
// 		p,
// 	}
// }
