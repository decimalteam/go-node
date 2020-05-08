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
		Sell time.Duration `json:"sell"`
	} `json:"timeout_ms"`
	CountAccounts int `json:"count_accounts"`
}

func (c *Config) UnmarshalJSON(data []byte) (err error) {
	var tmp struct {
		TimeoutMs struct {
			Send string `json:"send"`
			Buy  string `json:"buy"`
			Sell string `json:"sell"`
		} `json:"timeout_ms"`
	}
	if err = json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	c.TimeoutMs = struct {
		Send time.Duration `json:"send"`
		Buy  time.Duration `json:"buy"`
		Sell time.Duration `json:"sell"`
	}{}

	c.TimeoutMs.Send, err = time.ParseDuration(tmp.TimeoutMs.Send)
	c.TimeoutMs.Buy, err = time.ParseDuration(tmp.TimeoutMs.Buy)
	c.TimeoutMs.Sell, err = time.ParseDuration(tmp.TimeoutMs.Sell)
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
