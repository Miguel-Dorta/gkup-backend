package main

func getInterface(args map[string]interface{}, k string) interface{} {
	v, ok := args[k]
	if !ok {
		printError(k + " not defined")
	}
	return v
}

func getInt(args map[string]interface{}, key string) int {
	i, ok := getInterface(args, key).(int)
	if !ok {
		printError("error parsing " + key)
	}
	return i
}

func getInt64(args map[string]interface{}, key string) int64 {
	i, ok := getInterface(args, key).(int64)
	if !ok {
		printError("error parsing " + key)
	}
	return i
}

func getString(args map[string]interface{}, key string) string {
	str, ok := getInterface(args, key).(string)
	if !ok {
		printError("error parsing " + key)
	}
	return str
}

func getStringSlice(args map[string]interface{}, key string) []string {
	strs, ok := getInterface(args, key).([]string)
	if !ok {
		printError("error parsing " + key)
	}
	return strs
}
