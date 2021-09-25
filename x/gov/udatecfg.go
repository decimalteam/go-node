package gov

import (
	"io/ioutil"
	"encoding/json"
)
type UpdateCFG struct {
	Url string 		`json:url`
	Address string `json:address`
}


func LoadUpgradeCFG(fileName string) *UpdateCFG {
	if !fileExist(fileName) {
		return nil
	}
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	updateCfg := new(UpdateCFG)
	
	err = json.Unmarshal(bytes, updateCfg)
	if err != nil {
		panic(err)
	}

	return updateCfg
}