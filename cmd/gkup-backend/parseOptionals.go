package main

func getOptionalInterface(args map[string]interface{}, k string) interface{} {
	v, ok := args[k]
	if !ok {
		return nil
	}
	return v
}

func getOptionalInt(args map[string]interface{}, key string) int {
	in := getOptionalInterface(args, key)
	if in == nil {
		return 0
	}
	i, ok := in.(int)
	if !ok {
		printError("error parsing " + key)
	}
	return i
}

func getOptionalString(args map[string]interface{}, key string) string {
	in := getOptionalInterface(args, key)
	if in == nil {
		return ""
	}
	s, ok := in.(string)
	if !ok {
		printError("error parsing " + key)
	}
	return s
}
