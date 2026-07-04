package main

import (
	"sync"
)

var (
	Handlers = map[string]func([]Value) Value{
		"PING":    ping,
		"SET":     set,
		"GET":     get,
		"HSET":    hset,
		"HGET":    hget,
		"HGETALL": hgetall,
	}

	HSETs   = map[string]map[string]string{}
	SETs    = map[string]string{}
	SETsMu  = sync.RWMutex{}
	HSETsMu = sync.RWMutex{}
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
	val := args[1].bulk

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

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "error: HSET requires 3 arguments."}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "error: HGET command requires 2 arguments only."}
	}

	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "error: HGETALL command requires 1 arguments."}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	innerMap, ok := HSETs[hash]
	if !ok {
		HSETsMu.RUnlock()
		return Value{typ: "null"}
	}

	results := make([]Value, 0, len(innerMap)*2)

	for k, v := range innerMap {
		results = append(results, Value{typ: "bulk", bulk: k})
		results = append(results, Value{typ: "bulk", bulk: v})
	}

	HSETsMu.RUnlock()

	return Value{typ: "array", array: results}
}
