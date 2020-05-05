package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Config struct {
	TimeoutMs struct {
		Send time.Duration `json:"send"`
		Buy  time.Duration `json:"buy"`
	} `json:"timeout_ms"`
}

func ImportConfig(path string) (Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.Unmarshal(file, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
