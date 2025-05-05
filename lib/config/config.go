package config

import (
	"CloudTusk/lib/log"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

type ConfigParams struct {
	SearchValue interface{}
	LastModTime time.Time
	mutex       sync.Mutex
}

var baseDir = "config/"

const extension = ".json"

func (cp *ConfigParams) Get(fileName, path string) *ConfigParams {
	pathOfKeys := strings.Split(path, "->")

	fileName += extension

	cp.readFile(baseDir+fileName, pathOfKeys)

	return cp
}

func (cp *ConfigParams) readFile(filepath string, pathOfKeys []string) (err error) {
	cp.mutex.Lock()

	fileStat, err := os.Stat(filepath)

	if err == nil {
		if fileStat.ModTime().After(cp.LastModTime) {
			file, _ := os.ReadFile(filepath)

			var mapData map[string]interface{}

			err = json.Unmarshal(file, &mapData)

			cp.SearchValue = valueByKey(mapData, pathOfKeys)
			cp.LastModTime = fileStat.ModTime()
		}
	}

	cp.mutex.Unlock()

	return err
}

func (cp *ConfigParams) String() string {
	return fmt.Sprintf("%v", cp.SearchValue)
}

func valueByKey(value map[string]interface{}, keys []string) interface{} {
	nestedValue := value[keys[0]]
	count := len(keys)

	for i := 1; i < count; i++ {
		nestedValue = nestedValue.(map[string]interface{})[keys[i]]
	}

	if nestedValue == nil {
		log.Error("value not found for key: " + reflect.ValueOf(value).String())
	}

	return nestedValue
}
