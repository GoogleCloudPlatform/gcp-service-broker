// Copyright 2018 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// BrokerResponse encodes an OSB HTTP response in a (technical) human and
// machine readable way.
type BrokerResponse struct {
	// WARNING: BrokerResponse is exposed to users and automated tooling
	// so DO NOT remove or rename fields unless strictly necessary.
	// You MAY add new fields.
	Error        error           `json:"error,omitempty"`
	Url          string          `json:"url,omitempty"`
	Method       string          `json:"http_method,omitempty"`
	StatusCode   int             `json:"status_code,omitempty"`
	ResponseBody json.RawMessage `json:"response,omitempty"`
}

func (br *BrokerResponse) UpdateError(err error) {
	if br.Error == nil {
		br.Error = err
	}
}

func (br *BrokerResponse) UpdateRequest(req *http.Request) {
	if req == nil {
		return
	}

	br.Url = req.URL.String()
	br.Method = req.Method
}

func (br *BrokerResponse) UpdateResponse(res *http.Response) {
	if res == nil {
		return
	}

	br.StatusCode = res.StatusCode

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		br.UpdateError(err)
	} else {
		br.ResponseBody = json.RawMessage(body)
	}
}

func (br *BrokerResponse) InError() bool {
	return br.Error != nil
}

func (br *BrokerResponse) String() string {
	if br.InError() {
		return fmt.Sprintf("%s %s -> %d, Error: %q)", br.Method, br.Url, br.StatusCode, br.Error)
	}

	return fmt.Sprintf("%s %s -> %d, %q", br.Method, br.Url, br.StatusCode, br.ResponseBody)
}
