package testutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// GetJSONString extracts the string from the root JSON that's located at the
// given path
// E.g. for {"hello": {"world": "val"}} looking up "hello.world" would return
// "val"
func GetJSONString(root *json.RawMessage, path string) (string, error) {
	if strings.Contains(path, ".") {
		split := strings.Split(path, ".")
		subobj, err := GetNestedJSON(root, split[0])
		if err != nil {
			return "", err
		}
		return GetJSONString(subobj, strings.Join(split[1:], "."))
	}
	var m map[string]json.RawMessage
	err := json.Unmarshal(*root, &m)
	if err != nil {
		return "", err
	}
	var s string
	err = json.Unmarshal(m[path], &s)
	if err != nil {
		return "", err
	}
	return s, nil
}

// GetNestedJSON extracts a subobject from the given root, identified by the
// given key
func GetNestedJSON(root *json.RawMessage, key string) (*json.RawMessage,
	error) {
	var subobj map[string]*json.RawMessage
	err := json.Unmarshal(*root, &subobj)
	if err != nil {
		return nil, err
	}
	return subobj[key], nil
}

// ExtractPayload pulls out a raw JSON message from the given request's body.
func ExtractPayload(req *http.Request) (*json.RawMessage, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var root json.RawMessage
	return &root, json.Unmarshal(body, &root)
}
