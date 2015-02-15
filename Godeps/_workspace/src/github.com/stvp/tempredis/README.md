tempredis
=========

Tempredis makes it easy to start and stop temporary `redis-server`
processes with custom configs for testing.

[API documentation](http://godoc.org/github.com/stvp/tempredis)

Example
-------

```go
package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/stvp/tempredis"
)

func main() {
	server, err := tempredis.Start(
		tempredis.Config{
			"port":      "11001",
			"databases": "8",
		},
	)
	if err != nil {
		panic(err)
	}
	defer server.Term()

	conn, err := redis.Dial("tcp", ":11001")
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	conn.Do("SET", "foo", "bar")
}
```

Or, even easier:

```go
package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/stvp/tempredis"
)

func main() {
  config := tempredis.Config{
    "port":      "11001",
    "databases": "8",
  }

  tempredis.Temp(config, func(err error) {
    if err != nil {
      panic(err)
    }

    conn, err := redis.Dial("tcp", ":11001")
    defer conn.Close()
    if err != nil {
      panic(err)
    }

    conn.Do("SET", "foo", "bar")
  })
}
```

If you don't care about normal shutdown behavior or want to simulate a crash,
you can send a KILL signal to the server with:

```go
server.Kill()
```

Should I use this outside of testing?
-------------------------------------

No.

