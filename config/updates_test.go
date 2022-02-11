package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadFromOldConfig(t *testing.T) {
	oldconfig := []byte(`{"last_update":8123401,
	"all_updates":{"https://repo.decimalchain.com/7010801":7010801,
	"https://repo.decimalchain.com/7098701":7098701,
	"https://repo.decimalchain.com/7348701":7348701,
	"https://repo.decimalchain.com/7519401":7519401,
	"https://repo.decimalchain.com/7944001":7944001,
	"https://repo.decimalchain.com/7980901":7980901,
	"https://repo.decimalchain.com/8037701":8037701,
	"https://repo.decimalchain.com/8123401":8123401}}`)
	updinf := NewUpdatesInfo("")
	r := bytes.NewReader(oldconfig)
	updinf.load(r)
	assert.Equal(t, 8, len(updinf.AllBlocks), "old AllBlock must be 8")
	assert.Equal(t, int64(8123401), updinf.LastBlock, "LastBlock must be in safe")
}

func TestSaveLoad(t *testing.T) {
	updinf := NewUpdatesInfo("")
	updinf.PushNewPlanHeight(1)
	updinf.AddExecutedPlan("1", 1)
	updinf.AddExecutedPlan("2", 1)
	updinf.AddExecutedPlan("3", 1)
	tmp := make([]byte, 0)
	buf := bytes.NewBuffer(tmp)
	updinf.save(buf)
	newinf := NewUpdatesInfo("")
	newinf.load(buf)
	assert.Equal(t, updinf.AllBlocks, newinf.AllBlocks, "AllBlocks must be same")
	assert.Equal(t, updinf.LastBlock, newinf.LastBlock, "LastBlock must be same")
}
