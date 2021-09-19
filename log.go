package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func initializeLogger() error {
	var file io.WriteCloser
	// we use StdOut
	file = os.Stdout
	if len(gConfig.LogFile) > 0 && gConfig.LogFile != "-" {
		var err error
		if file, err = os.OpenFile(gConfig.LogFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0700); err != nil {
			return fmt.Errorf("Cannot create log file %q: %v", gConfig.LogFile, err)
		}
	}
	gLogger = log.New(file, "", 0)
	return nil
}

func accessLoggerFunc(r *http.Request, err error) {
	if len(gConfig.LogFile) == 0 {
		return
	}
	user := r.Context().Value("user").(string)
	t := r.Context().Value("time").(time.Time)
	timestr := t.Format("Jan 02 2006 15:04:05.000")
	if err == nil {
		err = fmt.Errorf("")
	}
	switch r.Method {
	case "COPY", "MOVE":
		dst := ""
		if u, err := url.Parse(r.Header.Get("Destination")); err == nil {
			dst = u.Path
		}
		o := r.Header.Get("Overwrite")
		gLogger.Printf("%s %s %s %s -> %s o=%s %v", timestr, user, r.Method, r.URL.Path, dst, o, err)
	default:
		gLogger.Printf("%s %s %s %s %v", timestr, user, r.Method, r.URL.Path, err)
	}
}
