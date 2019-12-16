package errors

import (
	"encoding/json"
	"fmt"
)

type errorStruct struct {
	Code int32  `json:"code"`
	Text string `json:"text"`
}

func (e *errorStruct) String() string {
	return fmt.Sprintf("error: (%d) %s", e.Code, e.Text)
}

func (e *errorStruct) Error() string {
	return e.Text
}

func (e *errorStruct) Marshal() []byte {
	data, _ := json.Marshal(e)
	return data
}

func New(text string, code ...int32) error {
	if len(code) > 0 {
		return &errorStruct{code[0], text}
	}
	return &errorStruct{0, text}
}
