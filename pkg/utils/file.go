package utils

import "strings"

type File string

// Filename returns name of file
func (fp File) Filename() string {
	if fp == "" {
		return ""
	}

	file := strings.Split(string(fp), "/")
	last := file[len(file)-1]
	return last
}

// Path returns path of file
func (fp File) Path() string {
	if fp == "" {
		return ""
	}

	file := strings.Split(string(fp), "/")
	if len(file) == 1 {
		return file[0]
	}

	return strings.Join(file[:len(file)-2], "/")
}
