package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	dayFlag  = flag.Int("day", 0, "Day number to initialize.")
	initFlag = flag.Bool("init", false, "Initialize a new day.")
)

var functionRegistry = map[int]func(){}

func main() {
	flag.Parse()

	day := *dayFlag
	if day == 0 {
		fmt.Printf("--day is required.\n")
		return
	}

	if f, ok := functionRegistry[day]; ok {
		f()
	} else if *initFlag {
		write := func(path string, data string) {
			if _, err := os.Stat(path); err == nil {
				fmt.Printf("%s already exists.\n", path)
				return
			}
			if err := os.WriteFile(path, []byte(data), 0644); err != nil {
				panic(err)
			}
			fmt.Printf("wrote %s.\n", path)
		}
		write(fmt.Sprintf("data/day%d.txt", day), "")
		write(fmt.Sprintf("day%d.go", day), fmt.Sprintf(dayTemplate, day))
	} else {
		fmt.Printf("No function registered for day %d. Use --init to create one.\n", day)
	}
}

const dayTemplate = `package main

import "fmt"

func day%[1]d() {
	for _, line := range readDataLines("day%[1]d.txt") {
		fmt.Printf("%%s\n", line)
	}
}

func init() {
	functionRegistry[%[1]d] = day%[1]d
}
`
