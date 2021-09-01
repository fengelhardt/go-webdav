package main

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/hjson/hjson-go"
)

func pathClean(p string) string {
	endswithslash := strings.HasSuffix(p, "/")
	r := path.Clean(p)
	if endswithslash {
		r = fmt.Sprintf("%s/", r)
	}
	return r
}

func unmarshalHJson(data []byte, v interface{}) error {
	var hjdat interface{}
	if err := hjson.Unmarshal(data, &hjdat); err != nil {
		return err
	}
	jdat, _ := json.Marshal(hjdat)
	if err := json.Unmarshal(jdat, v); err != nil {
		return err
	}
	return nil
}
