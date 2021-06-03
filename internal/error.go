package internal

import (
	"fmt"
)

// ThrowError 抛出错误
func ThrowError(err error) {
	if err != nil {
		panic(err)
	}
}

// ThrowErrorMessage 抛出错误
func ThrowErrorMessage(template string, args ...interface{}) {
	if len(args) == 0 {
		panic(template)
	}
	msg := fmt.Sprintf(template, args...)
	panic(msg)
}