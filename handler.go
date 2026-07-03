package main

import (
	"sync"
)

var (
	Handlers = map[string]func([]Value) Value{
		"PING": ping,
		"SET":  set,
		"GET":  get,
	}

	SETs   = map[string]string{}
	SETsMu = sync.RWMutex{}
)

func ping(args []Value) Value {
	if len(args) == 0 {

		return Value{typ: "string", str: "PONG"}
	}
	return Value{typ: "string", str: args[0].bulk}
}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "error: wrong no. of arguments for SET command. Requires 2 only."}
	}

	key := args[0].bulk
	val := args[0].bulk

	SETsMu.Lock()
	SETs[key] = val
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "error: GET command requires 1 argument only."}
	}

	key := args[0].bulk

	SETsMu.RLock()
	val, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: val}
}
