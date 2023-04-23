package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNode_RunJoinNetwork(t *testing.T) {

	//c04d6b84fd4fc978
	//node := new(Node)
	//err := node.RunJoinNetwork("c04d6b84fd4fc978")
	//if err != nil {
	//	t.Fail()
	//}

	destUrl := fmt.Sprintf("%s/api/v1/users/user/%s/network/%s/join", userUrl, "1", "dfc82f28aa6dcebc")
	resp, err := http.Post(destUrl, "application/json", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
