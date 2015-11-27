package kh

import (
	"fmt"
	"strings"
)

type Response struct {
	Stdout  []string
	Stderr  []string
	Log     []string
	Verbose bool
}

func (r *Response) SetVerbose(verbose bool) {
	r.Verbose = verbose
}

func (r *Response) Debugf(format string, args ...interface{}) {
	if r.Verbose {
		r.WriteLog(fmt.Sprintf(format, args))
	}
}

func (r *Response) Debug(msg string) {
	if r.Verbose {
		r.WriteLog(msg)
	}
}

func (r *Response) SprintLog() string {
	return strings.Join(r.Log, "\n")
}

func (r *Response) SprintStdout() string {
	return strings.Join(r.Stdout, "\n")
}

func (r *Response) SprintStderr() string {
	return strings.Join(r.Stderr, "\n")
}

func (r *Response) WriteLog(entry string) {
	r.Log = append(r.Log, entry)
}

func (r *Response) WriteStdout(entry string) {
	r.Stdout = append(r.Stdout, entry)
}

func (r *Response) WriteStderr(entry string) {
	r.Stderr = append(r.Stderr, entry)
}
