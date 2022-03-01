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

type BitData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type BitDataNull struct {
	ID   *string `json:"id,omitempty"`
	Name *string `json:"name"`
}

type NameData struct {
	Name string `json:"name"`
}

type IdData struct {
	ID string `json:"id"`
}
