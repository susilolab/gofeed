package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"reflect"
)

type Object interface {
	toJson() ([]byte, error)
}

// Enkripsi password dengan sha1
func hashPassword(password string) string {
	sha := sha1.New()
	sha.Write([]byte(password))
	enc := sha.Sum(nil)
	return fmt.Sprintf("%x", enc)
}

func (ul *UserLogin) toJson() ([]byte, error) {
	res, err := json.Marshal(ul)
	return res, err
}

func (u *User) toJson() ([]byte, error) {
	res, err := json.Marshal(u)
	return res, err
}

func isSet(a interface{}, key interface{}) bool {
	av := reflect.ValueOf(a)
	kv := reflect.ValueOf(key)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		if int64(av.Len()) > kv.Int() {
			return true
		}
	case reflect.Map:
		if kv.Type() == av.Type().Key() {
			return av.MapIndex(kv).IsValid()
		}
	}

	return false
}
