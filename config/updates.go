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

func (plan *updatesInfo) Push(name string, height int64) error {
	plan.LastBlock = height
	plan.AllBlocks[name] = height

	bytes, err := json.Marshal(plan)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(plan.filename, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (plan *updatesInfo) Load() error {
	if !fileExist(plan.filename) {
		err := ioutil.WriteFile(plan.filename, []byte("{}"), 0644)
		if err != nil {
			return err
		}
	}

	bytes, err := ioutil.ReadFile(plan.filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, plan)
	if err != nil {
		return err
	}

	return nil
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
