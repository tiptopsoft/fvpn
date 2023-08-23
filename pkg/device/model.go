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

package device

type Instance struct {
	UserId string `json:"userId"`
	Addr   string `json:"addr"`
	Status string `json:""` //live/died
}

type Response struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"data"`
	Message string      `json:"msg"`
}

type JoinRequest struct {
	NetWorkId string `json:"networkId"`
	CIDR      string `json:"Cidr"`
}

type JoinResponse struct {
	CIDR    string `json:"Cidr"`
	IP      string `json:"ip"`
	Name    string `json:"name"`
	Network string `json:"network"`
}

type LeaveRequest struct {
	JoinRequest
}

type LeaveResponse struct {
	JoinResponse
}

type InitResponse struct {
	IP    string `json:"ip"`
	Mask  string `json:"mask"`
	AppId string `json:"appId"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	UserId string `json:"userId"`
}

func HttpOK(data interface{}) Response {
	return Response{
		Code:    200,
		Result:  data,
		Message: "",
	}
}

func HttpError(message string) Response {
	return Response{
		Code:    500,
		Result:  nil,
		Message: message,
	}
}
