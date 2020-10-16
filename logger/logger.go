package logger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

type Interface interface {
	Deubg(...interface{})
	Info(...interface{})
}

type Logger struct {
	LevelDebug bool
	LevelInfo  bool
}

func New() Logger { return Logger{} }

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
				dump, e := httputil.DumpResponse(t, true)
				if e != nil {
					l.Debug("Dump Response Error:", e)
				}
				l.Debug("Response:", string(dump))

			default:
				log.Println(t)
			}
		}
	}
}

func (l Logger) Info(messages ...interface{}) {
	if l.LevelInfo {
		for _, m := range messages {
			fmt.Println(m)
		}
	}
}
