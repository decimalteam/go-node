package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type updatesInfo struct {
	filename  string
	LastBlock int64            `json:"last_update"`
	AllBlocks map[string]int64 `json:"all_updates"`
}

func NewUpdatesInfo(planfile string) *updatesInfo {
	return &updatesInfo{
		filename:  planfile,
		LastBlock: 0,
		AllBlocks: make(map[string]int64),
	}
}

func (plan *updatesInfo) Push(name string, height int64) {
	plan.LastBlock = height
	plan.AllBlocks[name] = height

	bytes, err := json.Marshal(plan)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(plan.filename, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func (plan *updatesInfo) Load() *updatesInfo {
	if !fileExist(plan.filename) {
		err := ioutil.WriteFile(plan.filename, []byte("{}"), 0600)
		if err != nil {
			panic(err)
		}
	}

	bytes, err := ioutil.ReadFile(plan.filename)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytes, plan)
	if err != nil {
		panic(err)
	}

	return plan
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
