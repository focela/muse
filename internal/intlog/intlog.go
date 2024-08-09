// Copyright (c) 2024 Focela Technologies. All rights reserved.
// Use of this source code is governed by an MIT style
// license that can be found in the LICENSE file.

// Package intlog provides internal logging for Aid's internal development usage only.
package intlog

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/focela/aid/debug"
	"github.com/focela/aid/internal/utils"
)

const (
	stackFilterKey = "/internal/intlog"
)

// Print prints `v` with a newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Print(ctx context.Context, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprint(v...), false)
}

// Printf prints `v` with format `format` using fmt.Printf.
// The parameter `v` can be multiple variables.
func Printf(ctx context.Context, format string, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprintf(format, v...), false)
}

// Error prints `v` with a newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Error(ctx context.Context, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprint(v...), true)
}

// Errorf prints `v` with format `format` using fmt.Printf.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	if !utils.IsDebugEnabled() {
		return
	}
	doPrint(ctx, fmt.Sprintf(format, v...), true)
}

// PrintFunc prints the output from function `f`.
// It only calls function `f` if debug mode is enabled.
func PrintFunc(ctx context.Context, f func() string) {
	if !utils.IsDebugEnabled() {
		return
	}
	s := f()
	if s == "" {
		return
	}
	doPrint(ctx, s, false)
}

// ErrorFunc prints the output from function `f`.
// It only calls function `f` if debug mode is enabled.
func ErrorFunc(ctx context.Context, f func() string) {
	if !utils.IsDebugEnabled() {
		return
	}
	s := f()
	if s == "" {
		return
	}
	doPrint(ctx, s, true)
}

// doPrint handles the actual logging with optional stack trace printing.
func doPrint(ctx context.Context, content string, stack bool) {
	if !utils.IsDebugEnabled() {
		return
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05.000"))
	buffer.WriteString(" [INTE] ")
	buffer.WriteString(file())
	buffer.WriteString(" ")
	if s := traceIdStr(ctx); s != "" {
		buffer.WriteString(s + " ")
	}
	buffer.WriteString(content)
	buffer.WriteString("\n")
	if stack {
		buffer.WriteString("Caller Stack:\n")
		buffer.WriteString(debug.StackWithFilter(stackFilterKey))
	}
	fmt.Print(buffer.String())
}

// traceIdStr retrieves and returns the trace ID string for logging output.
func traceIdStr(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if traceId := spanCtx.TraceID(); traceId.IsValid() {
		return "{" + traceId.String() + "}"
	}
	return ""
}

// file returns the caller file name along with its line number.
func file() string {
	_, p, l := debug.CallerWithFilter([]string{stackFilterKey})
	return fmt.Sprintf(`%s:%d`, filepath.Base(p), l)
}
