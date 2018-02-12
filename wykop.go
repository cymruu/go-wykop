package wykop

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const _useragent = "go-wykop 0.1.0"
const _BASEURL = "https://a.wykop.pl"

type WykopAPI struct {
	baseURL       string
	APIkey        string
	secretKey     string
	connectionKey string
	login         string
	userKey       string
	httpClient    http.Client
	errorHanders  map[uint16]func(*ErrorResponse, *WykopRequest)
}

func Create(APIkey, secretKey string) *WykopAPI {
	return &WykopAPI{APIkey: APIkey, secretKey: secretKey, baseURL: _BASEURL, httpClient: http.Client{Timeout: 10 * time.Second}, errorHanders: make(map[uint16]func(*ErrorResponse, *WykopRequest))}
}
func (c *WykopAPI) AddErrorHandler(errorCode uint16, f func(*ErrorResponse, *WykopRequest)) {
	c.errorHanders[errorCode] = f
}
func decodeJSON(data []byte, target interface{}) error {
	return json.Unmarshal(data, target)
}
func (c *WykopAPI) handleWykopError(wykopError *ErrorResponse, request *WykopRequest) {
	fmt.Printf("Handling error: %v", wykopError)
	if handler, ok := c.errorHanders[wykopError.ErrorObject.Code]; ok == true {
		handler(wykopError, request)
	}
}
func (c *WykopAPI) signRequest(request *WykopRequest) {
	checkBase := c.secretKey + request.buildURL()
	if request.postData != nil {
		var keys []string
		for k := range request.postData {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, v := range keys {
			checkBase += request.postData[v][0] + ","
		}
		if len(checkBase) > 0 {
			checkBase = checkBase[:len(checkBase)-1]
		}
	}
	checksum := md5.Sum([]byte(checkBase))
	request.checksum = fmt.Sprintf("%x", checksum)
}
func (c *WykopAPI) NewRequest(endpoint string, options ...RequestOptional) *WykopRequest {
	var APIParams APIParamsT
	if c.APIkey != "" {
		APIParams = append(APIParams, APIParamPair{"appkey", c.APIkey})
	}
	if c.userKey != "" {
		APIParams = append(APIParams, APIParamPair{"userkey", c.userKey})
	}
	options = append(options, OptionAPIParams(APIParams))
	return newRequest(endpoint, options...)
}
func (c *WykopAPI) sendRequest(request *WykopRequest, target interface{}) error {
	requestMethod := request.method()
	c.signRequest(request)
	req, _ := http.NewRequest(requestMethod, request.buildURL(), strings.NewReader(request.postData.Encode()))
	switch requestMethod {
	case "POST":
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("apisign", request.checksum)
	req.Header.Add("User-Agent", _useragent)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Print(err)
		return err
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	APIErr := ErrorResponse{}
	decodeJSON(data, &APIErr)
	if APIErr.ErrorObject.Code > 0 == true {
		c.handleWykopError(&APIErr, request)
		return &APIErr
	}
	err = decodeJSON(data, target)
	if err != nil {
		return err
	}
	return nil
}

func (c *WykopAPI) Login(login, connectionKey string) bool {
	c.login = login
	c.connectionKey = connectionKey
	postData := url.Values{}
	postData.Add("login", c.login)
	postData.Add("accountkey", c.connectionKey)
	req := c.NewRequest("user/login", OptionPostData(postData))
	resp := AuthorizationResponse{}
	err := c.sendRequest(req, &resp)
	if err != nil {
		return false
	}
	c.userKey = resp.Userkey
	return true
}

func (c *WykopAPI) GetEntry(entryID string) *EntryResponse {
	resp := EntryResponse{}
	err := c.sendRequest(c.NewRequest("entries/index", OptionMethodParams(MethodParamsT{entryID})), &resp)
	if err != nil {
		return nil
	}
	return &resp
}

func (c *WykopAPI) GetConversationList() *[]ConversationListItem {
	resp := []ConversationListItem{}
	err := c.sendRequest(c.NewRequest("pm/conversationslist"), &resp)
	if err != nil {
		return nil
	}
	return &resp
}

func (c *WykopAPI) Observe(username string) bool {
	req := c.NewRequest("profile/observe", OptionMethodParams(MethodParamsT{username}))
	var resp string
	err := c.sendRequest(req, &resp)
	if err != nil || resp != WYKOP_TRUE_RESPONSE {
		return false
	}
	return true
}
func (c *WykopAPI) Unobserve(username string) bool {
	req := c.NewRequest("profile/unobserve", OptionMethodParams(MethodParamsT{username}))
	var resp string
	err := c.sendRequest(req, &resp)
	if err != nil || resp != WYKOP_TRUE_RESPONSE {
		return false
	}
	return true
}
