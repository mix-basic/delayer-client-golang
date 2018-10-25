package delayer

import (
	"github.com/gomodule/redigo/redis"
	"errors"
	"time"
)

// 键名
const (
	KEY_JOP_POOL       = "delayer:jop_pool"
	PREFIX_JOP_BUCKET  = "delayer:jop_bucket:"
	PREFIX_READY_QUEUE = "delayer:ready_queue:"
)

// 客户端结构
type Client struct {
	Conn     redis.Conn
	Host     string
	Port     string
	Database int
	Password string
}

// 初始化
func (p *Client) Init() error {
	// 创建连接
	if p.Conn == nil {
		conn, err := redis.Dial("tcp", p.Host+":"+p.Port)
		if err != nil {
			return err
		}
		p.Conn = conn
	}
	// 验证密码
	if (p.Password != "") {
		if _, err := p.Conn.Do("AUTH", p.Password); err != nil {
			p.Conn.Close()
			return err
		}
	}
	// 选库
	if _, err := p.Conn.Do("SELECT", p.Database); err != nil {
		p.Conn.Close()
		return err
	}
	return nil
}

// 增加任务
func (p *Client) Push(message Message, delayTime int, readyMaxLifetime int) (bool, error) {
	// 参数验证
	if !message.Valid() {
		return false, errors.New("Invalid message.");
	}
	// 增加
	p.Conn.Send("MULTI")
	p.Conn.Send("HMSET", PREFIX_JOP_BUCKET+message.ID, "topic", message.Topic, "body", message.Body)
	p.Conn.Send("EXPIRE", PREFIX_JOP_BUCKET+message.ID, delayTime+readyMaxLifetime)
	p.Conn.Send("ZADD", KEY_JOP_POOL, time.Now().Unix()+int64(delayTime), message.ID)
	_, err := p.Conn.Do("EXEC")
	if err != nil {
		return false, err
	}
	// 返回
	return true, nil
}

// 取出任务
func (p *Client) Pop(topic string) (*Message, error) {
	id, err := redis.String(p.Conn.Do("RPOP", PREFIX_READY_QUEUE+topic))
	if err != nil || id == "" {
		return nil, err
	}
	result, err := redis.StringMap(p.Conn.Do("HGETALL", PREFIX_JOP_BUCKET+id))
	if err != nil {
		return nil, err
	}
	if result["topic"] == "" || result["body"] == "" {
		return nil, errors.New("")
	}
	err = p.Conn.Send("DEL", PREFIX_JOP_BUCKET+id)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		ID:    id,
		Topic: result["topic"],
		Body:  result["body"],
	}
	return msg, nil
}

// 阻塞取出任务
func (p *Client) BPop(topic string, timeout int) (*Message, error) {
	id, err := redis.String(p.Conn.Do("BRPOP", PREFIX_READY_QUEUE+topic, timeout))
	if err != nil || id == "" {
		return nil, err
	}
	result, err := redis.StringMap(p.Conn.Do("HGETALL", PREFIX_JOP_BUCKET+id))
	if err != nil {
		return nil, err
	}
	if result["topic"] == "" || result["body"] == "" {
		return nil, errors.New("")
	}
	err = p.Conn.Send("DEL", PREFIX_JOP_BUCKET+id)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		ID:    id,
		Topic: result["topic"],
		Body:  result["body"],
	}
	return msg, nil
}

// 移除任务
func (p *Client) Remove(id string) (bool, error) {
	p.Conn.Send("MULTI")
	p.Conn.Send("ZREM", KEY_JOP_POOL, id)
	p.Conn.Send("DEL", PREFIX_JOP_BUCKET+id)
	values, err := redis.Values(p.Conn.Do("EXEC"))
	if err != nil {
		return false, err
	}
	v := values[0].(int64)
	v1 := values[1].(int64)
	if v == 0 || v1 == 0 {
		return false, nil
	}
	// 返回
	return true, nil
}
