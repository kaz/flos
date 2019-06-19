package state

import (
	"encoding/json"
	"io/ioutil"

	"github.com/mattn/go-jsonpointer"
)

func Get(path string) (interface{}, error) {
	mu.RLock()
	defer mu.RUnlock()
	return jsonpointer.Get(store, path)
}
func Put(path string, data interface{}) error {
	mu.Lock()
	defer mu.Unlock()

	err := jsonpointer.Set(store, path, data)
	if err != nil {
		return err
	}

	rawStore, err = json.Marshal(store)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(STORE_FILE, rawStore, 0644)
}
