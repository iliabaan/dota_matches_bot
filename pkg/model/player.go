package model

type Player struct {
	Profile struct {
		ID   int    `json:"id"`
		Name string `json:"personaname"`
	} `json:"profile"`
}
