package model

type Expression struct {
	ID   int    `json:"id"`
	Exp  string `json:"exp"`
	Resp float64
}
