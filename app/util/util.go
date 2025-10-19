package util

import "grimoire/handlers"

func UpdateArgs(args *handlers.Arg, key string, value any) *handlers.Arg {
	newArgs := *args
	newArgs[key] = value
	
	return &newArgs
}