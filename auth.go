package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
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

func hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

func getNonce() string {
	return "abcdef"
}

func unquote(s string) string {
	return strings.Trim(s, "\"")
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

func parsePropertyString(s string) map[string]string {
	propmap := make(map[string]string)
	props := strings.Split(s, ",")
	for _, p := range props {
		if p != "" {
			kv := strings.Split(p, "=")
			propmap[strings.ToLower(trim(kv[0]))] = trim(unquote(kv[1]))
		}
	}
	return propmap
}

func checkPropertyMap(prop map[string]string) bool {
	ret := true
	var ok bool
	_, ok = prop["username"]
	ret = ret && ok
	_, ok = prop["realm"]
	ret = ret && ok
	_, ok = prop["uri"]
	ret = ret && ok
	_, ok = prop["nonce"]
	ret = ret && ok
	_, ok = prop["response"]
	ret = ret && ok
	return ret
}
