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
	Limits        *wykopLimits
}

func Create(APIkey, secretKey string) *WykopAPI {
	return &WykopAPI{APIkey: APIkey, secretKey: secretKey, baseURL: _BASEURL, httpClient: http.Client{Timeout: 10 * time.Second}, errorHanders: make(map[uint16]func(*ErrorResponse, *WykopRequest))}
}
func (c *WykopAPI) InitializeLimits(options ...limitOptional) {
	c.Limits = initializeLimits()
	for _, op := range options {
		op(c.Limits)
	}
}
func (c *WykopAPI) AddErrorHandler(errorCode uint16, f func(*ErrorResponse, *WykopRequest)) {
	c.errorHanders[errorCode] = f
}
func decodeJSON(data []byte, target interface{}) error {
	switch v := target.(type) {
	case *string:
		*v = string(data)
		return nil
	default:
		return json.Unmarshal(data, target)
	}
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
	if c.Limits != nil {
		c.Limits.register()
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	APIErr := ErrorResponse{}
	decodeJSON(data, &APIErr)
	if APIErr.ErrorObject.Code > 0 {
		c.handleWykopError(&APIErr, request)
		return &APIErr
	}
	err = decodeJSON(data, target)
	if err != nil {
		APIErr.ErrorObject.Code = errMalformedReponse
		APIErr.ErrorObject.Message = "Returned data is neither APIErr nor target object. Make sure the target argument is valid"
		return &APIErr
	}
	return nil
}

func (c *WykopAPI) Login(login, connectionKey string) (bool, error) {
	c.login = login
	c.connectionKey = connectionKey
	postData := url.Values{}
	postData.Add("login", c.login)
	postData.Add("accountkey", c.connectionKey)
	req := c.NewRequest("user/login", OptionPostData(postData))
	resp := AuthorizationResponse{}
	if err := c.sendRequest(req, &resp); err != nil {
		return false, err
	}
	c.userKey = resp.Userkey
	return true, nil
}

func (c *WykopAPI) GetEntry(entryID string) (*EntryResponse, error) {
	resp := EntryResponse{}
	err := c.sendRequest(c.NewRequest("entries/index", OptionMethodParams(MethodParamsT{entryID})), &resp)
	return &resp, err
}

func (c *WykopAPI) GetConversationList() (*[]ConversationListItem, error) {
	resp := []ConversationListItem{}
	err := c.sendRequest(c.NewRequest("pm/conversationslist"), &resp)
	return &resp, err
}
func (c *WykopAPI) observeHandler(endpoint string, username string) (bool, error) {
	req := c.NewRequest(fmt.Sprintf("profile/%s", endpoint), OptionMethodParams(MethodParamsT{username}))
	var resp string
	err := c.sendRequest(req, &resp)
	if err != nil || resp != WYKOP_TRUE_RESPONSE {
		return false, err
	}
	return true, nil
}
func (c *WykopAPI) Observe(username string) (bool, error) {
	return c.observeHandler("observe", username)
}
func (c *WykopAPI) Unobserve(username string) (bool, error) {
	return c.observeHandler("unobserve", username)
}
func (c *WykopAPI) GetNotifications(page uint) (*[]Notification, error) {
	APIParams := APIParamsT{APIParamPair{"page", fmt.Sprint(page)}}
	var resp []Notification
	err := c.sendRequest(c.NewRequest("mywykop/notifications", OptionAPIParams(APIParams)), &resp)
	return &resp, err
}
