package app

import "strings"

func newTask(file string) *task {
	fileparts := strings.Split(file, "/")

	return &task{
		filename: fileparts[len(fileparts)-1],
		path:     file,
	}
}

type task struct {
	filename string
	path     string
	err      error
}
