package scanner

import (
	"bufio"
	"os"
	"strings"
	"todoapp/todo"
)

type Scanner struct {
	todoList   *todo.List
	sliceEvent []Event
}

func NewScanner(todoList *todo.List) *Scanner {
	return &Scanner{
		todoList: todoList,
	}
}

func (s *Scanner) Start() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printPrompt()

		ok := scanner.Scan()

		if !ok {
			return
		}

		inputString := scanner.Text()

		result := s.process(inputString)

		if result != "" {
			if result == needExit {
				printExit()
				return
			}

			printResult(result)
		}

		event := NewEvent(result, inputString)

		s.sliceEvent = append(s.sliceEvent, event)

	}
}

func (s *Scanner) process(inputString string) string {
	fields := strings.Fields(inputString)

	if len(fields) == 0 {
		return emptyInput
	}

	cmd := fields[0]

	if cmd == "exit" {
		return needExit
	}

	if cmd == "add" {
		return s.cmdAdd(fields)
	}

	if cmd == "list" {
		return s.cmdList(fields)
	}

	if cmd == "done" {
		return s.cmdDone(fields)
	}

	if cmd == "del" {
		return s.cmdDel(fields)
	}

	if cmd == "help" {
		return s.cmdHelp(fields)
	}

	if cmd == "events" {
		return s.cmdEvents(fields)
	}

	return unknownCommand

}

func (s *Scanner) cmdAdd(fields []string) string {
	if len(fields) < 3 {
		return wrongArgs
	}

	title := fields[1]

	taskText := ""

	for i := 2; i < len(fields); i++ {
		if i == len(fields)-1 {
			taskText += fields[i]
		}

		taskText += fields[i] + " "
	}

	task := todo.NewTask(title, taskText)

	s.todoList.AddTask(task)

	printAdd(title)

	return ""
}

func (s *Scanner) cmdList(fields []string) string {
	if len(fields) != 1 {
		return wrongArgs
	}

	tasks := s.todoList.ListTasks()

	printTasks(tasks)

	return ""

}

func (s *Scanner) cmdDone(fields []string) string {
	if len(fields) != 2 {
		return wrongArgs
	}

	title := fields[1]

	_, err := s.todoList.CompleteTask(title)

	if err != nil {
		return err.Error()
	}

	printDone(title)

	return ""
}

func (s *Scanner) cmdDel(fields []string) string {
	if len(fields) != 2 {
		return wrongArgs
	}

	title := fields[1]

	deleteTaskResult := s.todoList.DeleteTask(title)

	if deleteTaskResult != nil {
		return deleteTaskResult.Error()
	}

	printDelete(title)

	return ""

}

func (s *Scanner) cmdHelp(fields []string) string {
	if len(fields) != 1 {
		return wrongArgs
	}

	printHelp()

	return ""
}

func (s *Scanner) cmdEvents(fields []string) string {
	if len(fields) != 1 {
		return wrongArgs
	}

	printEvents(s.sliceEvent)

	return ""

}
