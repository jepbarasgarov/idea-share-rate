package responses

import (
	"errors"
)

type GeneralResponse struct {
	Success bool        `json:"success"`
	ErrMsg  string      `json:"err_msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var ErrOK = errors.New("OK")

const (
	PG_CODE_UNIQUE_VIOLATION string = "23505"
)

type Language struct {
	TM string `json:"TM"`
	RU string `json:"RU"`
	EN string `json:"EN"`
}
type Lang string

const (
	Turkmen Lang = "TM"
	Russian Lang = "RU"
	English Lang = "EN"
)
