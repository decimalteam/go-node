package gov

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type SkipPlan struct {
	filename string
	mapping  map[string]int64
}

func NewSkipPlan(planfile string) *SkipPlan {
	return &SkipPlan{
		filename: planfile,
		mapping:  make(map[string]int64),
	}
}

func (plan *SkipPlan) Push(name string, height int64) {
	skipPlans := plan.Load()
	skipPlans[name] = height

	bytes, err := json.Marshal(plan.mapping)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(plan.filename, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

func (plan *SkipPlan) Load() map[string]int64 {
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

	err = json.Unmarshal(bytes, &plan.mapping)
	if err != nil {
		panic(err)
	}

	return plan.mapping
}

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
