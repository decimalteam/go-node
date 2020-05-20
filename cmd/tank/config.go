package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Config struct {
	Timeout       map[string]time.Duration `json:"timeout"`
	Weights       map[string]int           `json:"weights"`
	CountAccounts int                      `json:"count_accounts"`
}

func (c *Config) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		Timeout       map[string]string `json:"timeout"`
		Weights       map[string]int    `json:"weights"`
		CountAccounts int               `json:"count_accounts"`
	}
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	c.Timeout = make(map[string]time.Duration)
	for tx, timeout := range tmp.Timeout {
		c.Timeout[tx], err = time.ParseDuration(timeout)
		if err != nil {
			return err
		}
	}

	c.CountAccounts = tmp.CountAccounts
	c.Weights = tmp.Weights
	return err
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
