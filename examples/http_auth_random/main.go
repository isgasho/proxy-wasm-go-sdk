// Copyright 2020 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"hash/fnv"
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const clusterName = "httpbin"

func main() {
	proxywasm.SetNewHttpContext(newContext)
}

type httpHeaders struct {
	// you must embed the default context so that you need not to reimplement all the methods by yourself
	proxywasm.DefaultContext
	contextID uint32
}

func newContext(contextID uint32) proxywasm.HttpContext {
	return &httpHeaders{contextID: contextID}
}

// override default
func (ctx *httpHeaders) OnHttpRequestHeaders(int, bool) types.Action {
	hs, err := proxywasm.HostCallGetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCritical("failed to get request headers: ", err.Error())
		return types.ActionContinue
	}
	for _, h := range hs {
		proxywasm.LogInfo("request header: ", h[0], ": ", h[1])
	}

	if _, err := proxywasm.HostCallDispatchHttpCall(
		clusterName, hs, "", [][2]string{}, 50000); err != nil {
		proxywasm.LogCritical("dipatch httpcall failed: ", err.Error())
	}

	proxywasm.LogInfo("http call dispatched to ", clusterName)

	return types.ActionPause
}

// override default
func (ctx *httpHeaders) OnHttpCallResponse(_ int, bodySize int, _ int) {
	hs, err := proxywasm.HostCallGetHttpCallResponseHeaders()
	if err != nil {
		proxywasm.LogCritical("failed to get response body: ", err.Error())
		return
	}

	for _, h := range hs {
		proxywasm.LogInfo("response header from httpbin: ", h[0], ": ", h[1])
	}

	b, err := proxywasm.HostCallGetHttpCallResponseBody(0, bodySize)
	if err != nil {
		proxywasm.LogCritical("failed to get response body: ", err.Error())
		proxywasm.HostCallResumeHttpRequest()
		return
	}

	s := fnv.New32a()
	if _, err := s.Write(b); err != nil {
		proxywasm.LogCritical("failed to calculate hash: ", err.Error())
		proxywasm.HostCallResumeHttpRequest()
		return
	}

	if s.Sum32()%2 == 0 {
		proxywasm.LogInfo("access granted")
		proxywasm.HostCallResumeHttpRequest()
		return
	}

	msg := "access forbidden"
	proxywasm.LogInfo(msg)
	proxywasm.HostCallSendHttpResponse(403, [][2]string{
		{"powered-by", "proxy-wasm-go-sdk!!"},
	}, msg)
}

// override default
func (ctx *httpHeaders) OnLog() {
	proxywasm.LogInfo(strconv.FormatUint(uint64(ctx.contextID), 10), " finished")
}
