package panics

import (
	"fmt"
	"regexp"
	"runtime/debug"
)

var panicRegex = regexp.MustCompile(`runtime/panic\.go[ \S]+\s[ \S]+\s+([\S]+\.go:\d+)`)

type Error struct {
	message string
	reason  interface{}
	stack   []byte
}

func NewError(message string, reason interface{}) *Error {
	return &Error{
		message: message,
		reason:  reason,
		stack:   debug.Stack(),
	}
}

func (r *Error) Error() string {
	result := panicRegex.FindSubmatch(r.stack)

	if len(result) > 1 {
		return fmt.Sprintf("%s panic: %s in %s", r.message, r.reason, string(result[1]))
	}

	return fmt.Sprintf("%s panic: %s stack: %s", r.message, r.reason, string(r.stack))
}

func (r *Error) StackTrace() string {
	return string(r.stack)
}
