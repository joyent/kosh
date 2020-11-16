/*
Package logger is a simple logging module based upon the advice of Dave Cheney
	* https://dave.cheney.net/2015/11/05/lets-talk-about-logging
	* https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern
*/
package logger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// Interface is the default interface for logging with this logger. Debug() is
// for developer targeted output. While Info is for (verbose) user targeted
// output.
type Interface interface {
	Debug(...interface{})
	Info(...interface{})
}

type NullLogger struct{}

func (nl NullLogger) Debug(msgs ...interface{}) {}
func (nl NullLogger) Info(msgs ...interface{})  {}

// Logger is the default logger with configuration levels for debug (developer)
// output, and info (verbose user) output.
type Logger struct {
	LevelDebug bool
	LevelInfo  bool
}

// New returns a new instance of the Logger struct
func New() Logger { return Logger{} }

// Debug outputs developer targeted messaging.
func (l Logger) Debug(messages ...interface{}) {
	if l.LevelDebug {
		for _, m := range messages {
			if m == nil {
				continue
			}
			switch t := m.(type) {
			case *http.Request:
				dump, e := httputil.DumpRequestOut(t, true)
				if e != nil {
					log.Println("Got error:", e)
				}
				log.Println("Request:", string(dump))
			case *http.Response:
				dump, e := httputil.DumpResponse(t, false)
				if e != nil {
					l.Debug(fmt.Sprintf("Dump Response Error: %s", e))
				}
				l.Debug("Response:", string(dump))

			default:
				log.Println(t)
			}
		}
	}
}

// Info outputs more verbose user targed information
func (l Logger) Info(messages ...interface{}) {
	if l.LevelInfo {
		for _, m := range messages {
			fmt.Println(m)
		}
	}
}
