package main

// (c) 2021 Frank Engelhardt, <frank@f9e.de>

import (
	"fmt"
	"os"
)

type UserConfig map[string]struct {
	Hash  string `json:"Hash"`
	Quota string `json:"Quota"`
}

func loadUserConfig(file string) (UserConfig, error) {
	var ret UserConfig
	data, err := os.ReadFile(file)
	if err == nil {
		err = unmarshalHJson(data, &ret)
	}
	if err != nil {
		return ret, fmt.Errorf("Error with user config %q: %s", file, err)
	} else {
		return ret, nil
	}
}

func (uc UserConfig) hasUser(u string) bool {
	_, ok := uc[u]
	return ok
}

func (uc UserConfig) checkHash(user, hash string) bool {
	h := uc[user]
	return hash == h.Hash
}
