package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func getArgs() (map[string]interface{}, error) {
	d := json.NewDecoder(os.Stdin)

	var args map[string]interface{}
	if err := d.Decode(args); err != nil {
		return nil, fmt.Errorf("error decoding json: %s", err)
	}
	return args, nil
}

func getInt64(args map[string]interface{}, key string) int64 {
	i, ok := getInterface(args, key).(int64)
	if !ok {
		printError("error parsing " + key)
	}
	return i
}

func getStringSlice(args map[string]interface{}, key string) []string {
	strs, ok := getInterface(args, key).([]string)
	if !ok {
		printError("error parsing " + key)
	}
	return strs
}

func getString(args map[string]interface{}, key string) string {
	str, ok := getInterface(args, key).(string)
	if !ok {
		printError("error parsing " + key)
	}
	return str
}

func getInterface(args map[string]interface{}, k string) interface{} {
	v, ok := args[k]
	if !ok {
		printError(k + " not defined")
	}
	return v
}
