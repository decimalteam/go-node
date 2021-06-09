package errors

import "encoding/json"

type Err struct {
	Description string            `json:"description"`
	Params      map[string]string `json:"params,omitempty"`
}

type Param struct {
	Key   string
	Value string
}

func NewParam(key, value string) Param {
	return Param{
		Key:   key,
		Value: value,
	}
}

func Encode(description string, params ...Param) string {
	err := Err{
		Description: description,
	}

	if params != nil {
		err.Params = make(map[string]string)
		for _, param := range params {
			err.Params[param.Key] = param.Value
		}
	}

	result, _ := json.Marshal(err)
	return string(result)
}
