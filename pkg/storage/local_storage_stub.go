// +build !js !wasm

package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// LocalStorage provides a file-based storage for non-WebAssembly builds
type LocalStorage struct {
	dataDir string
}

func NewLocalStorage() *LocalStorage {
	// Use a local data directory
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".island-merge")
	os.MkdirAll(dataDir, 0755)
	
	return &LocalStorage{
		dataDir: dataDir,
	}
}

// Set stores a value in a local file
func (ls *LocalStorage) Set(key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	
	filePath := filepath.Join(ls.dataDir, key+".json")
	return os.WriteFile(filePath, jsonData, 0644)
}

// Get retrieves a value from a local file
func (ls *LocalStorage) Get(key string, target interface{}) error {
	filePath := filepath.Join(ls.dataDir, key+".json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	
	return json.Unmarshal(data, target)
}

// Remove deletes a key file
func (ls *LocalStorage) Remove(key string) {
	filePath := filepath.Join(ls.dataDir, key+".json")
	os.Remove(filePath)
}

// Exists checks if a key file exists
func (ls *LocalStorage) Exists(key string) bool {
	filePath := filepath.Join(ls.dataDir, key+".json")
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// Clear removes all files in the data directory
func (ls *LocalStorage) Clear() {
	os.RemoveAll(ls.dataDir)
	os.MkdirAll(ls.dataDir, 0755)
}

// GetKeys returns all keys that match a prefix
func (ls *LocalStorage) GetKeys(prefix string) []string {
	var keys []string
	
	files, err := filepath.Glob(filepath.Join(ls.dataDir, prefix+"*.json"))
	if err != nil {
		return keys
	}
	
	for _, file := range files {
		base := filepath.Base(file)
		key := base[:len(base)-5] // Remove .json extension
		keys = append(keys, key)
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