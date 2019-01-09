package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const defaultConfigPath = "../conf/db_config.json"

var config map[string]interface{}

func LoadConfigs(configPath string) (err error) {
	if configPath == "" {
		configPath = defaultConfigPath
	}
	file, err := os.Open(configPath)
	if err != nil {
		return
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	fmt.Println("config=", config)
	return
}

func String(key string, defaultValue string) (value string) {
	value, _ = getString(key, defaultValue)
	return
}

func getString(key string, defaultValue string) (value string, found bool) {
	v, found := config[key]
	if !found {
		return defaultValue, false
	}
	if value, ok := v.(string); ok {
		return value, ok
	} else {
		return defaultValue, false
	}
}

func StringOrError(key string) (value string, err error) {
	value, found := getString(key, "")
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func StringOrPanic(key string) string {
	v, err := StringOrError(key)
	if err != nil {
		panic(err)
	}
	return v
}

func Int(key string, defaultValue int) (value int, found bool) {
	v, found := config[key]
	if !found {
		value = defaultValue
		return
	}

	if v64, found := v.(float64); found {
		return int(v64), found
	} else {
		return defaultValue, false
	}
}

func IntOrError(key string) (value int, err error) {
	value, found := Int(key, 0)
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func IntOrPanic(key string) int {
	v, err := IntOrError(key)
	if err != nil {
		panic(err)
	}
	return v
}

func Int64(key string, defaultValue int64) (value int64, found bool) {
	v, found := config[key]
	if !found {
		value = defaultValue
		return
	}

	if v64, found := v.(float64); found {
		return int64(v64), found
	} else {
		return defaultValue, false
	}
}

func Int64OrError(key string) (value int64, err error) {
	value, found := Int64(key, 0)
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func Int64OrPanic(key string) int64 {
	v, err := Int64OrError(key)
	if err != nil {
		panic(err)
	}
	return v
}

func Bool(key string, defaultValue bool) (value bool, found bool) {
	v, found := config[key]
	if !found {
		value = defaultValue
		return
	}

	if b, ok := v.(bool); ok {
		return b, ok
	} else if i, ok := v.(int); ok {
		if i == 1 {
			return true, true
		} else if i == 0 {
			return false, true
		} else {
			return defaultValue, false
		}
	} else {
		return defaultValue, false
	}
}

func BoolOrError(key string) (value bool, err error) {
	value, found := Bool(key, false)
	if !found {
		err = fmt.Errorf("%s is not configured", key)
	}
	return
}

func BoolOrPanic(key string) bool {
	v, err := BoolOrError(key)
	if err != nil {
		panic(err)
	}
	return v
}
