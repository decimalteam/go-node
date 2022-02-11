package config

import (
	"encoding/json"
	"io"
	"os"
)

type updatesInfo struct {
	filename   string
	PlanBlocks []int64          `json:"plan_blocks"` // all last plan blocks to calculate grace periods
	LastBlock  int64            `json:"last_update"` // last height of 'software_upgrade'
	AllBlocks  map[string]int64 `json:"all_updates"` // map of executed upgrades. key - plan name, value - height
}

func NewUpdatesInfo(planfile string) *updatesInfo {
	return &updatesInfo{
		filename:   planfile,
		PlanBlocks: make([]int64, 0),
		LastBlock:  0,
		AllBlocks:  make(map[string]int64),
	}
}

func (plan *updatesInfo) PushNewPlanHeight(planHeight int64, horizonHeight int64) {
	if planHeight > plan.LastBlock {
		plan.LastBlock = planHeight
	}
	doadd := true
	//do not add existing height
	for _, h := range plan.PlanBlocks {
		if h == planHeight {
			doadd = false
		}
	}
	if doadd {
		plan.PlanBlocks = append(plan.PlanBlocks, planHeight)
	}

	// cleanup all below horizon
	newblocks := make([]int64, 0, len(plan.PlanBlocks))
	for _, h := range plan.PlanBlocks {
		if h >= horizonHeight {
			newblocks = append(newblocks, h)
		}
	}
	plan.PlanBlocks = newblocks
}

func (plan *updatesInfo) SaveExecutedPlan(planName string, planHeight int64) {
	plan.AllBlocks[planName] = planHeight
}

func (plan *updatesInfo) Save() error {
	f, err := os.OpenFile(plan.filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return plan.save(f)
}

func (plan *updatesInfo) save(wrt io.Writer) error {
	bytes, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	_, err = wrt.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (plan *updatesInfo) Load() error {
	if !fileExist(plan.filename) {
		err := os.WriteFile(plan.filename, []byte("{}"), 0644)
		if err != nil {
			return err
		}
	}
	f, err := os.OpenFile(plan.filename, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return plan.load(f)
}

func (plan *updatesInfo) load(rdr io.Reader) error {
	bytes, err := io.ReadAll(rdr)
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
