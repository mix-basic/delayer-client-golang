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

## DEMO

### `push` 方法

放入一个任务。

```go
```

### `pop` 方法

取出一个到期的任务。

```go
```

### `bPop` 方法

阻塞取出一个到期的任务。

```go
```

### `remove` 方法

移除一个未到期的任务。

```go
```

## License

Apache License Version 2.0, http://www.apache.org/licenses/
