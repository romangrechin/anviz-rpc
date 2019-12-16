package models

type Response struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}
