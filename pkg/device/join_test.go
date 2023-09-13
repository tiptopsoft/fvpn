// Copyright 2023 TiptopSoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package device

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNode_RunJoinNetwork(t *testing.T) {

	//c04d6b84fd4fc978
	//node := new(Peer)
	//err := node.RunJoinNetwork("c04d6b84fd4fc978")
	//if err != nil {
	//	t.Fail()
	//}

	destUrl := fmt.Sprintf("%s/api/v1/users/user/%s/network/%s/join", "", "1", "dfc82f28aa6dcebc")
	resp, err := http.Post(destUrl, "application/json", nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
