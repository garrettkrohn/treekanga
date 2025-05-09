package cmd

import (
	"fmt"
	"reflect"
	"strings"

	com "github.com/garrettkrohn/treekanga/common"
)

func PrintConfig(config com.AddConfig) {
	printStruct(reflect.ValueOf(config), 0)
}

func printStruct(v reflect.Value, indent int) {
	t := v.Type()

	if v.Kind() == reflect.Struct {
		fmt.Printf("%s%s: {\n", getIndent(indent), t.Name())
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i)

			if value.Kind() == reflect.Struct {
				printStruct(value, indent+1)
			} else if value.Kind() == reflect.Ptr {
				if !value.IsNil() {
					fmt.Printf("%s%s: %v\n", getIndent(indent+1), field.Name, value.Elem())
				} else {
					fmt.Printf("%s%s: nil\n", getIndent(indent+1), field.Name)
				}
			} else {
				fmt.Printf("%s%s: %v\n", getIndent(indent+1), field.Name, value)
			}
		}
		fmt.Printf("%s}\n", getIndent(indent))
	} else {
		fmt.Println("Provided value is not a struct")
	}
}

func getIndent(indent int) string {
	return strings.Repeat("  ", indent)
}
