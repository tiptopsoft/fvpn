package model

type Instance struct {
	UserId string `json:"userId"`
	Addr   string `json:"addr"`
	Status string `json:""` //live/died
}
