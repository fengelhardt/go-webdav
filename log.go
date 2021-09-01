package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func accessLoggerFunc(r *http.Request, err error) {
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
		log.Printf("%s %s %s %s -> %s o=%s %v", timestr, user, r.Method, r.URL.Path, dst, o, err)
	default:
		log.Printf("%s %s %s %s %v", timestr, user, r.Method, r.URL.Path, err)
	}
}
