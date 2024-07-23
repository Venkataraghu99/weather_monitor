package models

import "encoding/json"

type AppError struct {
	ErrorMessage string `json:"errorMessage"`
}

func (err *AppError) Marshall() []byte {
	resErr, _ := json.Marshal(err)
	return resErr
}
