// Copyright (c) 2024 Focela Technologies. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rdebug

import (
	"bytes"
	"fmt"
	"runtime"
)

// PrintStack prints to standard error the stack trace returned by runtime.Stack.
func PrintStack(skip ...int) {
	fmt.Print(Stack(skip...))
}

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack(skip ...int) string {
	return StackWithFilter(nil, skip...)
}

// StackWithFilter returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
//
// The parameter `filter` is used to filter the path of the caller.
func StackWithFilter(filters []string, skip ...int) string {
	return StackWithFilters(filters, skip...)
}

// StackWithFilters returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
//
// The parameter `filters` is a slice of strings, which are used to filter the path of the
// caller.
//
// TODO Improve the performance using debug.Stack.
func StackWithFilters(filters []string, skip ...int) string {
	number := 0
	if len(skip) > 0 {
		number = skip[0]
	}
	var (
		name                  string
		space                 = "  "
		index                 = 1
		buffer                = bytes.NewBuffer(nil)
		ok                    = true
		pc, file, line, start = callerFromIndex(filters)
	)
	for i := start + number; i < maxCallerDepth; i++ {
		if i != start {
			pc, file, line, ok = runtime.Caller(i)
		}
		if ok {
			if filterFileByFilters(file, filters) {
				continue
			}
			if fn := runtime.FuncForPC(pc); fn == nil {
				name = "unknown"
			} else {
				name = fn.Name()
			}
			if index > 9 {
				space = " "
			}
			buffer.WriteString(fmt.Sprintf("%d.%s%s\n    %s:%d\n", index, space, name, file, line))
			index++
		} else {
			break
		}
	}
	return buffer.String()
}
