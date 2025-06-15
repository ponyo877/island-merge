// +build js,wasm

package storage

import (
	"encoding/json"
	"syscall/js"
)

// LocalStorage provides a Go interface to browser localStorage
type LocalStorage struct{}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

// Set stores a value in localStorage
func (ls *LocalStorage) Set(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	js.Global().Get("localStorage").Call("setItem", key, string(jsonData))
	return nil
}

// Get retrieves a value from localStorage
func (ls *LocalStorage) Get(key string, target interface{}) error {
	localStorage := js.Global().Get("localStorage")
	item := localStorage.Call("getItem", key)
	
	if item.IsNull() {
		return ErrNotFound
	}
	
	jsonStr := item.String()
	return json.Unmarshal([]byte(jsonStr), target)
}

// Remove deletes a key from localStorage
func (ls *LocalStorage) Remove(key string) {
	js.Global().Get("localStorage").Call("removeItem", key)
}

// Exists checks if a key exists in localStorage
func (ls *LocalStorage) Exists(key string) bool {
	localStorage := js.Global().Get("localStorage")
	item := localStorage.Call("getItem", key)
	return !item.IsNull()
}

// Clear removes all items from localStorage
func (ls *LocalStorage) Clear() {
	js.Global().Get("localStorage").Call("clear")
}

// GetKeys returns all keys in localStorage that match a prefix
func (ls *LocalStorage) GetKeys(prefix string) []string {
	localStorage := js.Global().Get("localStorage")
	length := localStorage.Get("length").Int()
	
	var keys []string
	for i := 0; i < length; i++ {
		key := localStorage.Call("key", i).String()
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			keys = append(keys, key)
		}
	}
	
	return keys
}

var ErrNotFound = &StorageError{"key not found"}

type StorageError struct {
	Message string
}

func (e *StorageError) Error() string {
	return e.Message
}