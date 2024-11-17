package dto

type Result struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"result"`
}
