package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/sha3"
)

func checkAuth(authheader, method string) (bool, string) {
	propstr := cutPropertyString(authheader)
	// check response sanity
	if propstr == "" {
		return false, ""
	}
	// check auth
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, []byte(propstr))
	hashstr := fmt.Sprintf("%x", hash)
	user := parseUser(propstr)
	if gUserConfig.hasUser(user) {
		if gUserConfig.checkHash(user, hashstr) {
			return true, user
		}
		return false, user
	}
	return false, user
}

func authenticateHeader() string {
	return fmt.Sprintf("Basic realm=%s", gConfig.Realm)
}

func parseUser(propstr string) string {
	decoded, _ := base64.StdEncoding.DecodeString(propstr)
	dstr := string(decoded)
	data := strings.Split(dstr, ":")
	if len(data) > 0 {
		return data[0]
	}
	return ""
}

func trim(s string) string {
	return strings.Trim(s, " \t\n\r")
}

func cutPropertyString(s string) string {
	if strings.HasPrefix(s, "Basic") || strings.HasPrefix(s, "basic") {
		return trim(s[5:])
	} else {
		return ""
	}
}
