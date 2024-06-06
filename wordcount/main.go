//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	res := make(map[string]int)
	for i := 0; i < len(args); i++ {
		file, err := os.Open(args[i])
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			key := scanner.Text()
			res[key]++
		}
	}
	for key, value := range res {
		if value >= 2 {
			fmt.Printf("%v\t%v\n", value, key)
		}
	}
}
