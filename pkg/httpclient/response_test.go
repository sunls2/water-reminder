package httpclient

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestResponse_PlainText(t *testing.T) {
	const text = "test response text."
	resp := newResponse([]byte(text))
	if resp.PlainText() != text {
		t.Errorf("plain text: (is) %s != (want) %s", resp.PlainText(), text)
	}
}

const (
	keyCode     = "code"
	code        = 100
	keyMsg      = "msg"
	msg         = "this is a message"
	keyNotExist = "not-exist"
)

var resp = newResponse([]byte(fmt.Sprintf(`{"%s": %d, "%s": "%s"}`, keyCode, code, keyMsg, msg)))

func TestResponse_GetString(t *testing.T) {
	s1, err := resp.GetString(keyMsg)
	if err != nil || s1 != msg {
		t.Errorf("GetString: (is) %s != (want) %s, error: %v", s1, msg, err)
	}
	_, err = resp.GetString(keyCode)
	if err == nil {
		t.Errorf("GetString: type mismatch should return error, but nil")
	}
}

func TestResponse_GetInt(t *testing.T) {
	c1, err := resp.GetInt(keyCode)
	if err != nil || c1 != code {
		t.Errorf("GetInt: (is) %d != (want) %d, error: %v", c1, code, err)
	}

	_, err = resp.GetInt(keyNotExist)
	if err == nil {
		t.Errorf("GetInt: when a key does not exist should return error, but nil")
	} else {
		t.Log(err)
	}

	_, err = resp.GetInt(keyMsg)
	if err == nil {
		t.Errorf("GetString: type mismatch should return error, but nil")
	}
}

const (
	keyMap   = "map"
	mStr     = `{"name": "ls", "age": 18}`
	keyArray = "array"
	arrayStr = `["lg", "dell", "aoc"]`
)

var (
	m     = make(map[string]any)
	array = make([]any, 0)
)

var resp2 = newResponse([]byte(fmt.Sprintf(`{"%s": %s, "%s": %s}`, keyMap, mStr, keyArray, arrayStr)))

func TestResponse_GetMap(t *testing.T) {
	if err := json.Unmarshal([]byte(mStr), &m); err != nil {
		t.Fatal(err)
	}
	m1, err := resp2.GetMap(keyMap)
	if err != nil {
		t.Errorf("GetMap: %v", err)
	}
	if equal := reflect.DeepEqual(m, m1); !equal {
		t.Errorf("GetMap: not equal: is %v want %v", m1, m)
	}
}

func TestResponse_GetArray(t *testing.T) {
	if err := json.Unmarshal([]byte(arrayStr), &array); err != nil {
		t.Fatal(err)
	}
	array1, err := resp2.GetArray(keyArray)
	if err != nil {
		t.Errorf("GetArray: %v", err)
	}
	if equal := reflect.DeepEqual(array, array1); !equal {
		t.Errorf("GetArray: not equal: is %v want %v", array1, array)
	}
}
