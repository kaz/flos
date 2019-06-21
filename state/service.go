package state

import (
	"encoding/json"

	"github.com/kaz/flos/camo"
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

	return camo.WriteFile(STORE_FILE, rawStore, 0644)
}
