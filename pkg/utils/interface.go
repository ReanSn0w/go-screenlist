package utils

type Logger interface {
	Logf(format string, args ...interface{})
}
