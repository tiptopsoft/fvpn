// Copyright 2023 Tiptopsoft, Inc.
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

package util

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestCheckNatType(t *testing.T) {

	s := "69 0 0 84 20 166 0 0 64 1 228 173 192 168 0 2 192 168 0 3 0 0 122 196 24 5 0 30 100 238 223 235 0 6 61 53 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55"
	arr := strings.Split(s, " ")
	var buff []byte
	for _, v := range arr {
		value, _ := strconv.Atoi(v)
		buff = append(buff, byte(value))
	}

	a := buff[0] >> 4
	fmt.Println(a)
}
