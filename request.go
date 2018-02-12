package wykop

import (
	"fmt"
	"net/url"
	"strings"
)

type MethodParamsT = []string
type APIParamPair struct {
	name  string
	value string
}
type APIParamsT = []APIParamPair
type RequestOptional func(*WykopRequest)
type WykopRequest struct {
	URL          string
	endpoint     string
	methodParams MethodParamsT
	APIParams    APIParamsT
	postData     url.Values
	checksum     string
}

func newRequest(endpoint string, options ...RequestOptional) *WykopRequest {
	req := &WykopRequest{}
	req.URL = _BASEURL
	req.endpoint = endpoint
	for _, op := range options {
		op(req)
	}
	return req
}
func OptionMethodParams(v MethodParamsT) RequestOptional {
	return func(r *WykopRequest) {
		r.methodParams = v
	}
}
func OptionAPIParams(v APIParamsT) RequestOptional {
	return func(r *WykopRequest) {
		for x := range v {
			r.APIParams = append(r.APIParams, v[x])
		}
	}
}
func OptionPostData(v url.Values) RequestOptional {
	return func(r *WykopRequest) {
		r.postData = v
	}
}
func (req *WykopRequest) method() string {
	if req.postData != nil {
		return "POST"
	}
	return "GET"
}
func (req *WykopRequest) buildURL() string {
	URL := fmt.Sprintf("%s/%s/", strings.TrimSuffix(req.URL, "/"), req.endpoint)
	if req.methodParams != nil {
		URL += fmt.Sprintf("%s/", strings.Join(req.methodParams, "/"))
	}
	if req.APIParams != nil {
		for x := range req.APIParams {
			URL += fmt.Sprintf("%s,%s,", req.APIParams[x].name, req.APIParams[x].value)
		}
	}

	return strings.TrimSuffix(URL, ",")
}
