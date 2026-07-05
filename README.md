# tiny-redis

A minimal Redis-compatible server implemented in Go. Supports basic key-value operations, hash maps, and persistence via AOF (Append Only File).

## Features

- **RESP Protocol**: Full implementation of Redis Serialization Protocol (RESP2)
- **Commands Supported**:
  - `PING` - Health check
  - `SET key value` - Set a key-value pair
  - `GET key` - Get a value by key
  - `HSET hash key value` - Set a field in a hash
  - `HGET hash key` - Get a field from a hash
  - `HGETALL hash` - Get all fields from a hash
- **Persistence**: AOF (Append Only File) for durability
- **Thread-Safe**: Uses `sync.RWMutex` for concurrent access
- **Default Port**: 6379 (Redis standard)

## Quick Start

### Prerequisites

- Go 1.21+

### Running the Server

```bash
go run main.go
```

Server starts listening on `:6379`

### Testing with redis-cli

```bash
# Ping
redis-cli ping
# PONG

# Set/Get
redis-cli set name "tiny-redis"
redis-cli get name

# Hash operations
redis-cli hset user:1 name "Alice" age "30"
redis-cli hget user:1 name
redis-cli hgetall user:1
```

## Project Structure

```
tiny-redis/
├── main.go      # Entry point, TCP server, connection handling
├── handler.go   # Command handlers (PING, SET, GET, HSET, HGET, HGETALL)
├── resp.go      # RESP protocol parser & serializer
├── aof.go       # Append-only file persistence
├── db.aof       # Persisted data (auto-generated)
├── go.mod       # Go module definition
└── .gitignore
```

## Architecture

### RESP Protocol (`resp.go`)

Parses and serializes Redis protocol types:
- **Simple Strings** (`+OK\r\n`)
- **Errors** (`-Error message\r\n`)
- **Integers** (`:1000\r\n`)
- **Bulk Strings** (`$6\r\nfoobar\r\n`)
- **Arrays** (`*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n`)
- **Null** (`$-1\r\n`)

### Command Handlers (`handler.go`)

All handlers follow the signature: `func(args []Value) Value`

Thread-safe using `sync.RWMutex`:
- `SETs` - String key-value store
- `HSETs` - Hash maps (map of maps)

### AOF Persistence (`aof.go`)

- Writes every `SET` and `HSET` command to `db.aof`
- Background goroutine syncs to disk every second
- On startup, replays AOF to restore state

## Example Session

```bash
$ redis-cli
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> SET greeting "Hello, tiny-redis!"
OK
127.0.0.1:6379> GET greeting
"Hello, tiny-redis!"
127.0.0.1:6379> HSET session:1 user "pabitra" role "admin"
OK
127.0.0.1:6379> HGET session:1 user
"pabitra"
127.0.0.1:6379> HGETALL session:1
1) "user"
2) "pabitra"
3) "role"
4) "admin"
```

## License

MIT
