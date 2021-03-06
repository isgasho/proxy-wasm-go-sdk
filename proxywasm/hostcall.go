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

package proxywasm

import (
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/rawhostcall"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

// wrappers on the rawhostcall package

func HostCallGetPluginConfiguration(dataSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypePluginConfiguration, 0, dataSize)
	return ret, types.StatusToError(st)
}

func HostCallGetVMConfiguration(dataSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeVMConfiguration, 0, dataSize)
	return ret, types.StatusToError(st)
}

func HostCallSendHttpResponse(statusCode uint32, headers [][2]string, body string) types.Status {
	shs := SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return rawhostcall.ProxySendLocalResponse(statusCode, nil, 0,
		stringBytePtr(body), len(body), hp, hl, -1,
	)
}

func hostCallSetEffectiveContext(contextID uint32) types.Status {
	return rawhostcall.ProxySetEffectiveContext(contextID)
}

func HostCallSetTickPeriodMilliSeconds(millSec uint32) error {
	return types.StatusToError(rawhostcall.ProxySetTickPeriodMilliseconds(millSec))
}

func HostCallGetCurrentTime() int64 {
	var t int64
	rawhostcall.ProxyGetCurrentTimeNanoseconds(&t)
	return t
}

func HostCallDispatchHttpCall(upstream string,
	headers [][2]string, body string, trailers [][2]string, timeoutMillisecond uint32) (uint32, error) {
	shs := SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)

	sts := SerializeMap(trailers)
	tp := &sts[0]
	tl := len(sts)

	var calloutID uint32

	u := stringBytePtr(upstream)
	switch st := rawhostcall.ProxyHttpCall(u, len(upstream),
		hp, hl, stringBytePtr(body), len(body), tp, tl, timeoutMillisecond, &calloutID); st {
	case types.StatusOK:
		currentState.registerCallout(calloutID)
		return calloutID, nil
	default:
		return 0, types.StatusToError(st)
	}
}

func HostCallGetHttpCallResponseHeaders() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpCallResponseHeaders)
	return ret, types.StatusToError(st)
}

func HostCallGetHttpCallResponseBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpCallResponseBody, start, maxSize)
	return ret, types.StatusToError(st)
}

func HostCallGetHttpCallResponseTrailers() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpCallResponseTrailers)
	return ret, types.StatusToError(st)
}

func HostCallDone() {
	switch st := rawhostcall.ProxyDone(); st {
	case types.StatusOK:
		return
	default:
		panic("unexpected status on proxy_done: " + strconv.FormatUint(uint64(st), 10))
	}
}

func HostCallGetDownStreamData(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeDownstreamData, start, maxSize)
	return ret, types.StatusToError(st)
}

func HostCallGetUpstreamData(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeUpstreamData, start, maxSize)
	return ret, types.StatusToError(st)
}

func HostCallGetHttpRequestHeaders() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpRequestHeaders)
	return ret, types.StatusToError(st)
}

func HostCallSetHttpRequestHeaders(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestHeaders, headers))
}

func HostCallGetHttpRequestHeader(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpRequestHeaders, key)
	return ret, types.StatusToError(st)
}

func HostCallRemoveHttpRequestHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestHeaders, key))
}

func HostCallSetHttpRequestHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

func HostCallAddHttpRequestHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestHeaders, key, value))
}

func HostCallGetHttpRequestBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpRequestBody, start, maxSize)
	return ret, types.StatusToError(st)
}

func HostCallGetHttpRequestTrailers() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpRequestTrailers)
	return ret, types.StatusToError(st)
}

func HostCallSetHttpRequestTrailers(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpRequestTrailers, headers))
}

func HostCallGetHttpRequestTrailer(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpRequestTrailers, key)
	return ret, types.StatusToError(st)
}

func HostCallRemoveHttpRequestTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpRequestTrailers, key))
}

func HostCallSetHttpRequestTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

func HostCallAddHttpRequestTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpRequestTrailers, key, value))
}

func HostCallResumeHttpRequest() error {
	return types.StatusToError(rawhostcall.ProxyContinueStream(types.StreamTypeRequest))
}

func HostCallGetHttpResponseHeaders() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpResponseHeaders)
	return ret, types.StatusToError(st)
}

func HostCallSetHttpResponseHeaders(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseHeaders, headers))
}

func HostCallGetHttpResponseHeader(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpResponseHeaders, key)
	return ret, types.StatusToError(st)
}

func HostCallRemoveHttpResponseHeader(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseHeaders, key))
}

func HostCallSetHttpResponseHeader(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

func HostCallAddHttpResponseHeader(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseHeaders, key, value))
}

func HostCallGetHttpResponseBody(start, maxSize int) ([]byte, error) {
	ret, st := getBuffer(types.BufferTypeHttpResponseBody, start, maxSize)
	return ret, types.StatusToError(st)
}

func HostCallGetHttpResponseTrailers() ([][2]string, error) {
	ret, st := getMap(types.MapTypeHttpResponseTrailers)
	return ret, types.StatusToError(st)
}

func HostCallSetHttpResponseTrailers(headers [][2]string) error {
	return types.StatusToError(setMap(types.MapTypeHttpResponseTrailers, headers))
}

func HostCallGetHttpResponseTrailer(key string) (string, error) {
	ret, st := getMapValue(types.MapTypeHttpResponseTrailers, key)
	return ret, types.StatusToError(st)
}

func HostCallRemoveHttpResponseTrailer(key string) error {
	return types.StatusToError(removeMapValue(types.MapTypeHttpResponseTrailers, key))
}

func HostCallSetHttpResponseTrailer(key, value string) error {
	return types.StatusToError(setMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

func HostCallAddHttpResponseTrailer(key, value string) error {
	return types.StatusToError(addMapValue(types.MapTypeHttpResponseTrailers, key, value))
}

func HostCallResumeHttpResponse() error {
	return types.StatusToError(rawhostcall.ProxyContinueStream(types.StreamTypeResponse))
}

func HostCallRegisterSharedQueue(name string) (uint32, error) {
	var queueID uint32
	ptr := stringBytePtr(name)
	st := rawhostcall.ProxyRegisterSharedQueue(ptr, len(name), &queueID)
	return queueID, types.StatusToError(st)
}

// TODO: not sure if the ABI is correct
func HostCallResolveSharedQueue(vmID, queueName string) (uint32, error) {
	var ret uint32
	st := rawhostcall.ProxyResolveSharedQueue(stringBytePtr(vmID),
		len(vmID), stringBytePtr(queueName), len(queueName), &ret)
	return ret, types.StatusToError(st)
}

func HostCallDequeueSharedQueue(queueID uint32) ([]byte, error) {
	var raw *byte
	var size int
	st := rawhostcall.ProxyDequeueSharedQueue(queueID, &raw, &size)
	if st != types.StatusOK {
		return nil, types.StatusToError(st)
	}
	return RawBytePtrToByteSlice(raw, size), nil
}

func HostCallEnqueueSharedQueue(queueID uint32, data []byte) error {
	return types.StatusToError(rawhostcall.ProxyEnqueueSharedQueue(queueID, &data[0], len(data)))
}

func HostCallGetSharedData(key string) (value []byte, cas uint32, err error) {
	var raw *byte
	var size int

	st := rawhostcall.ProxyGetSharedData(stringBytePtr(key), len(key), &raw, &size, &cas)
	if st != types.StatusOK {
		return nil, 0, types.StatusToError(st)
	}
	return RawBytePtrToByteSlice(raw, size), cas, nil
}

func HostCallSetSharedData(key string, data []byte, cas uint32) error {
	st := rawhostcall.ProxySetSharedData(stringBytePtr(key),
		len(key), &data[0], len(data), cas)
	return types.StatusToError(st)
}

func setMap(mapType types.MapType, headers [][2]string) types.Status {
	shs := SerializeMap(headers)
	hp := &shs[0]
	hl := len(shs)
	return rawhostcall.ProxySetHeaderMapPairs(mapType, hp, hl)
}

func getMapValue(mapType types.MapType, key string) (string, types.Status) {
	var rvs int
	var raw *byte
	if st := rawhostcall.ProxyGetHeaderMapValue(mapType, stringBytePtr(key), len(key), &raw, &rvs); st != types.StatusOK {
		return "", st
	}

	ret := RawBytePtrToString(raw, rvs)
	return ret, types.StatusOK
}

func removeMapValue(mapType types.MapType, key string) types.Status {
	return rawhostcall.ProxyRemoveHeaderMapValue(mapType, stringBytePtr(key), len(key))
}

func setMapValue(mapType types.MapType, key, value string) types.Status {
	return rawhostcall.ProxyReplaceHeaderMapValue(mapType, stringBytePtr(key), len(key), stringBytePtr(value), len(value))
}

func addMapValue(mapType types.MapType, key, value string) types.Status {
	return rawhostcall.ProxyAddHeaderMapValue(mapType, stringBytePtr(key), len(key), stringBytePtr(value), len(value))
}

func getMap(mapType types.MapType) ([][2]string, types.Status) {
	var rvs int
	var raw *byte

	st := rawhostcall.ProxyGetHeaderMapPairs(mapType, &raw, &rvs)
	if st != types.StatusOK {
		return nil, st
	}

	bs := RawBytePtrToByteSlice(raw, rvs)
	return DeserializeMap(bs), types.StatusOK
}

func getBuffer(bufType types.BufferType, start, maxSize int) ([]byte, types.Status) {
	var retData *byte
	var retSize int
	switch st := rawhostcall.ProxyGetBufferBytes(bufType, start, maxSize, &retData, &retSize); st {
	case types.StatusOK:
		// is this correct handling...?
		if retData == nil {
			return nil, types.StatusNotFound
		}
		return RawBytePtrToByteSlice(retData, retSize), st
	default:
		return nil, st
	}
}
