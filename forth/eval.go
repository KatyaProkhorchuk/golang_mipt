//go:build !solution

package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Stack []int

type Evaluator struct {
	// "стек", грамматика
	stack    Stack
	grammar  map[string]string
	function map[string]string
}

func (stack *Stack) Push(value int) {
	*stack = append(*stack, value)
}

func (stack *Stack) IsEmpty() bool {
	return len(*stack) == 0
}

func (stack *Stack) Top() int {
	return (*stack)[len(*stack)-1]
}

func (stack *Stack) Pop() int {
	if stack.IsEmpty() {
		return 0
	}
	top := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return top
}

func (stack *Stack) Size() int {
	return len(*stack)
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{stack: []int{}, grammar: map[string]string{}, function: map[string]string{
		"+":    "+",
		"-":    "-",
		"*":    "*",
		"/":    "/",
		"dup":  "dup",
		"over": "over",
		"drop": "drop",
		"swap": "swap",
	}}
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	commands := strings.Split(strings.ToLower(row), " ")
	var isErr error
	for len(commands) > 0 {
		command := commands[0]
		commands = commands[1:]
		if isErr != nil {
			break
		}
		if len(command) == 0 {
			continue
		}
		if cmd, ok := e.grammar[command]; ok {
			// fmt.Println(e.function[cmd], cmd, 1)
			if _, okk := e.function[cmd]; okk {
				command = cmd

			} else {
				_, isErr = e.Process(cmd)
				continue
			}
		}
		if val, err := strconv.Atoi(command); err == nil {
			// число -> добавляем в стек
			e.stack.Push(val)
		} else if command == "+" || command == "-" || command == "*" || command == "/" {
			isErr = e.Operation(command)
		} else if command == "dup" {
			isErr = e.Dup()
		} else if command == "drop" {
			isErr = e.Drop()
		} else if command == "over" {
			isErr = e.Over()
		} else if command == "swap" {
			isErr = e.Swap()
		} else if command == ":" {
			commands, isErr = e.AddNewCommand(commands)
		} else {
			isErr = fmt.Errorf("incorrect command")
		}
	}

	return e.stack, isErr
}

func (e *Evaluator) AddNewCommand(commands []string) ([]string, error) {
	if len(commands) == 0 {
		return commands, fmt.Errorf("incorrect new command")
	}
	command := commands[0]
	commands = commands[1:]
	if _, isErr := strconv.Atoi(command); isErr == nil {
		return commands, fmt.Errorf("command is number it's bad")
	}
	newcommand := ""
	for len(commands) > 0 && commands[0] != ";" {
		cmd := commands[0]
		if val, ok := e.grammar[cmd]; ok {
			cmd = val
		}
		newcommand += cmd + " "
		commands = commands[1:]
	}
	commands = commands[1:] // del ;
	if len(newcommand) == 0 {
		return commands, fmt.Errorf("empty new command")
	}
	e.grammar[command] = strings.TrimSpace(newcommand)
	return commands, nil
}

func (e *Evaluator) Operation(command string) error {
	// +, -, *, /
	if e.stack.Size() < 2 {
		return fmt.Errorf("size stack < 2")
	}
	if command == "+" {
		first := e.stack.Pop()
		second := e.stack.Pop()
		e.stack.Push(first + second)
	} else if command == "-" {
		first := e.stack.Pop()
		second := e.stack.Pop()
		e.stack.Push(second - first)
	} else if command == "*" {
		first := e.stack.Pop()
		second := e.stack.Pop()
		e.stack.Push(first * second)
	} else if command == "/" {
		first := e.stack.Pop()
		second := e.stack.Pop()
		if first == 0 {
			return fmt.Errorf("value is 0")
		}
		e.stack.Push(second / first)
	}
	return nil
}
func (e *Evaluator) Dup() error {
	if e.stack.Size() == 0 {
		return fmt.Errorf("invalid dup")
	}
	value := e.stack.Top()
	e.stack.Push(value)
	return nil
}
func (e *Evaluator) Over() error {
	if e.stack.Size() < 2 {
		return fmt.Errorf("invalid over")
	}
	first := e.stack.Pop()
	second := e.stack.Top()
	e.stack.Push(first)
	e.stack.Push(second)
	return nil
}
func (e *Evaluator) Drop() error {
	if e.stack.Size() == 0 {
		return fmt.Errorf("invalid drop")
	}
	e.stack.Pop()
	return nil
}
func (e *Evaluator) Swap() error {
	if e.stack.Size() < 2 {
		return fmt.Errorf("invalis swap")
	}
	first := e.stack.Pop()
	second := e.stack.Pop()
	e.stack.Push(first)
	e.stack.Push(second)
	return nil
}
