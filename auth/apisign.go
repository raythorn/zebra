package auth

import (
	"crypto/md5"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	APIMaxInterval int64 = 120
)

// Sort (url, timestamp="123", token=adad232, [salt=ada], [param1=value1], [param2=value2]...) in dictionary order,
// and join them into a string with ".", then sign the string use md5
//
//	Source data used to sign the api
//	source["url"] = "/user/add", MUST INCLUDE the request path without parameters
//	source["token"] = "ad234aadkjf-2=+", MUST INCLUDE the server signed token at login
//	source["timestamp"] = "1233132412", MUST INCLUDE the request timestamp
//	source["key"] = value, parameters attached to the request url, each saved as single data
//	source["salt"] = random data, provided by server, optional
type APISign struct {
	source map[string]string
}

func NewAPISign() *APISign {
	return &APISign{make(map[string]string)}
}

func (api *APISign) Set(key, value string) {
	api.source[key] = value
}

func (api *APISign) Get(key string) string {

	if value, ok := api.source[key]; ok {
		return value
	}

	return ""
}

func (api *APISign) Sign() (string, error) {

	if "" == api.Get("url") || "" == api.Get("token") || "" == api.Get("timestamp") {
		return "", errors.New("Sign: Token, URL and Timestamp cannot be empty.")
	}

	var sorted []string

	for key, value := range api.source {
		if value == "" {
			continue
		}

		var item string
		if key == "url" {
			item = value
		} else {
			item = key + "=" + value
		}
		sorted = append(sorted, item)
	}

	sort.Strings(sorted)

	signString := strings.Join(sorted, ".")
	hash := md5.New()
	_, err := hash.Write([]byte(signString))
	if nil != err {
		return "", errors.New("Sign: Hash failed")
	}

	signature := fmt.Sprintf("%x", hash.Sum(nil))

	return signature, nil
}

func (api *APISign) Verify(sign string) error {

	signature, err := api.Sign()
	if err != nil {
		return errors.New("Verify: Sign API failed")
	}

	if sign != signature {
		return errors.New("Verify: Signature not match")
	}

	ts := api.Get("timestamp")

	timestamp, err := strconv.ParseInt(ts, 10, 64)
	if nil != err {
		return errors.New("Verify: Parse timestamp failed")
	}

	now := time.Now().Unix()
	if interval := now - timestamp; interval > APIMaxInterval {
		return errors.New("Verify: API expired")
	}

	return nil
}
