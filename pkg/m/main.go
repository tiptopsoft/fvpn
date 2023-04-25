package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	type body struct {
		NetworkId string `json:"networkId"`
	}

	req := new(body)
	req.NetworkId = "xxxxxfdsa"

	buff, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buff))
}
