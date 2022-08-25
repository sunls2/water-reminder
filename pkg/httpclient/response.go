package httpclient

import (
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
)

const (
	keyPlainText = "plain-text"
)

type Response interface {
	PlainText() string

	GetString(key string) (string, error)
	GetInt(key string) (int, error)

	GetMap(key string) (map[string]any, error)
	GetArray(key string) ([]any, error)
}

type response map[string]any

func get[T any](resp response, key string) (T, error) {
	var value T
	v, ok := resp[key]
	if !ok {
		return value, errNotFoundKey(key)
	}
	if value, ok = v.(T); !ok {
		is := reflect.TypeOf(v)
		want := reflect.TypeOf(value)
		return value, errTypeMismatch(key, want.Name(), is.Name())
	}
	return value, nil
}

func (r response) GetMap(key string) (map[string]any, error) {
	return get[map[string]any](r, key)
}

func (r response) GetArray(key string) ([]any, error) {
	return get[[]any](r, key)
}

func (r response) GetString(key string) (string, error) {
	return get[string](r, key)
}

func (r response) GetInt(key string) (int, error) {
	f, err := get[float64](r, key)
	return int(f), err
}

func (r response) PlainText() string {
	v, _ := r[keyPlainText].(string)
	return v
}

var _ Response = (*response)(nil)

func newResponse(body []byte) Response {
	var response = response{keyPlainText: string(body)}
	// 尝试转换
	_ = json.Unmarshal(body, &response)
	return response
}

func errNotFoundKey(key string) error {
	return errors.Errorf("key '%s' not found", key)
}

func errTypeMismatch(key, want, is string) error {
	return errors.Errorf("key '%s' type mismatch: want %s is %s", key, want, is)
}
