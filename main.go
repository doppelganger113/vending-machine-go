package main

import (
	"flag"
	"fmt"
	"log"
	"vending-machine-go/internal"
)

var inputString, vendingMachineString string
var strict bool

func init() {
	flag.BoolVar(&strict, "strict", false, "strict input order")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Invalid number of arguments. Expecting 'input' 'vending_machine'")
	}
	inputString = args[0]
	vendingMachineString = args[1]
}

func getPattern() internal.PatternFunc {
	if strict == true {
		return internal.FindFirstPattern
	}
	return internal.FindFirstNoOrderPattern
}

// Usage:
// cmd products buckets
func main() {
	parsedInput, err := internal.ParseInput(inputString)
	if err != nil {
		fmt.Println(err)
		return
	}

	vendingMachine, err := internal.CreateFromString(vendingMachineString)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = internal.FindAndPopByOrder(vendingMachine, parsedInput, getPattern())
	if err != nil {
		fmt.Println(err)
		return
	}

	internal.PrintPretty(vendingMachine)
}
