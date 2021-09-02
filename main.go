package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/webdav"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	var err error
	gConfig, err = loadConfig(*gConfFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	gUserConfig, err = loadUserConfig(gConfig.UsersFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	h := &webdav.Handler{
		FileSystem: newUserDirFileSystem(),
		LockSystem: webdav.NewMemLS(),
		Logger:     accessLoggerFunc,
	}

	http.Handle(gConfig.URIPrefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), "time", time.Now()))
		authheader := r.Header.Get("Authorization")
		ok, user := checkAuth(authheader, r.Method)
		r = r.WithContext(context.WithValue(r.Context(), "user", user))
		if !ok {
			w.Header().Add("WWW-Authenticate", authenticateHeader())
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
			accessLoggerFunc(r, fmt.Errorf("Unauthorized"))
			return
		}
		h.ServeHTTP(w, r)
	}))

	addr := fmt.Sprintf(":%d", gConfig.Port)
	log.Printf("Serving %v%s", addr, gConfig.URIPrefix)
	log.Fatal(http.ListenAndServeTLS(addr, gConfig.TLSCertFile, gConfig.TLSKeyFile, nil))
	//log.Fatal(http.ListenAndServe(addr, nil))
}
