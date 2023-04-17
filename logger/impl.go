package logger

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
)

func Debug(arg any, args ...interface{}) {
	if message, ok := arg.(string); ok {
		logging("debug", true, message, args...)
	} else {
		args = []interface{}{arg}
		logging("debug", true, "%v", args...)
	}
}

func Info(message string, args ...interface{}) {
	logging("info", false, message, args...)
}

func Error(arg any, args ...interface{}) {
	if message, ok := arg.(string); ok {
		logging("error", true, message, args...)
	} else {
		args = []interface{}{arg}
		logging("error", true, "%v", args...)
	}
}

func getCallerInfo() (string, int) {
	_, fp, line, _ := runtime.Caller(3)
	_, fc, _, _ := runtime.Caller(0)
	path := strings.ReplaceAll(fp, filepath.Join(fc, "..", "..", ".."), "")

	return path, line
}

func logging(level string, withCallerInfo bool, message string, args ...interface{}) {
	body := fmt.Sprintf(message, args...)
	header := fmt.Sprintf("%s:", level)
	tail := ""
	if withCallerInfo {
		path, line := getCallerInfo()
		tail = fmt.Sprintf("[%s: %d]", path, line)
	}
	log.Printf("%s %s %s", header, body, tail)
}
