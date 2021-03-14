package main

import "errors"

var errKeyNotFound error = errors.New("key not found")

var storage = map[string]string{}

// Read reads a single key from storage
func Read(key string) (string, error) {
	val, ok := storage[key]
	if !ok {
		return "", errKeyNotFound
	}

	return val, nil
}

// Put puts a key, value pair in storage
func Put(key string, val string) bool {
	storage[key] = val
	_, ok := storage[key]
	return ok
}
