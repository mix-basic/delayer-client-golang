## Delayer Golang 客户端

客户端使用非常简单，提供了 `push`、`pop`、`bPop`、`remove` 四个方法操作任务。

## 安装

通过 `go get` 安装使用：

```shell
// install Redis client
go get github.com/gomodule/redigo/redis
// install Delayer client
go get github.com/mixstart/delayer-client-golang/delayer
```

## Documentation

- [API Reference](https://godoc.org/github.com/mixstart/delayer-client-golang/delayer)
- [Examples](https://godoc.org/github.com/mixstart/delayer-client-golang/delayer#pkg-examples)

## Example

### 创建客户端

通过连接信息创建客户端

```go
cli := delayer.Client{
    Host:     "127.0.0.1",
    Port:     "6379",
    Database: 0,
    Password: "",
}
cli.Init()
```

通过已有连接创建客户端

```go
pool := redis.Pool{}
conn := pool.Get();
cli := delayer.Client{
    Conn:conn,
}
cli.Init()
```

### `push` 方法

放入一个任务。

```go
cli := delayer.Client{
    Host:     "127.0.0.1",
    Port:     "6379",
    Database: 0,
    Password: "",
}
cli.Init()
msg := delayer.Message{
    ID:    fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()))),
    Topic: "test",
    Body:  "test body",
}
reply, err := cli.Push(msg, 10, 600)
fmt.Println(msg)
fmt.Println(reply)
fmt.Println(err)
```

### `pop` 方法

取出一个到期的任务。

```go
cli := delayer.Client{
    Host:     "127.0.0.1",
    Port:     "6379",
    Database: 0,
    Password: "",
}
cli.Init()
msg, err := cli.Pop("test");
fmt.Println(msg)
fmt.Println(err)
```

### `bPop` 方法

阻塞取出一个到期的任务。

```go
cli := delayer.Client{
    Host:     "127.0.0.1",
    Port:     "6379",
    Database: 0,
    Password: "",
}
cli.Init()
msg, err := cli.BPop("test", 10);
fmt.Println(msg)
fmt.Println(err)
```

### `remove` 方法

移除一个未到期的任务。

```go
cli := delayer.Client{
    Host:     "127.0.0.1",
    Port:     "6379",
    Database: 0,
    Password: "",
}
cli.Init()
ok, err := cli.Remove("9a8482a06630840ce7da9da62d748b8a")
fmt.Println(ok)
fmt.Println(err)
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
