package delayer

import (
	"github.com/gomodule/redigo/redis"
	"errors"
	"time"
)

// 键名
const (
	KEY_JOB_POOL       = "delayer:job_pool"
	PREFIX_JOB_BUCKET  = "delayer:job_bucket:"
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
	}
	return nil
}

// 增加任务
func (p *Client) Push(message Message, delayTime int, readyMaxLifetime int) (bool, error) {
	// 参数验证
	if !message.Valid() {
		return false, errors.New("Invalid message.");
	}
	// 执行事务
	p.Conn.Send("MULTI")
	p.Conn.Send("HMSET", PREFIX_JOB_BUCKET+message.ID, "topic", message.Topic, "body", message.Body)
	p.Conn.Send("EXPIRE", PREFIX_JOB_BUCKET+message.ID, delayTime+readyMaxLifetime)
	p.Conn.Send("ZADD", KEY_JOB_POOL, time.Now().Unix()+int64(delayTime), message.ID)
	values, err := redis.Values(p.Conn.Do("EXEC"))
	if err != nil {
		return false, err
	}
	// 事务结果处理
	v := values[0].(string)
	v1 := values[1].(int64)
	v2 := values[2].(int64)
	if v != "OK" || v1 == 0 || v2 == 0 {
		return false, nil
	}
	// 返回
	return true, nil
}

// 取出任务
func (p *Client) Pop(topic string) (*Message, error) {
	id, err := redis.String(p.Conn.Do("RPOP", PREFIX_READY_QUEUE+topic))
	if err != nil {
		return nil, err
	}
	result, err := redis.StringMap(p.Conn.Do("HGETALL", PREFIX_JOB_BUCKET+id))
	if err != nil {
		return nil, err
	}
	if result["topic"] == "" || result["body"] == "" {
		return nil, errors.New("Job bucket has expired or is incomplete")
	}
	p.Conn.Do("DEL", PREFIX_JOB_BUCKET+id)
	msg := &Message{
		ID:    id,
		Topic: result["topic"],
		Body:  result["body"],
	}
	return msg, nil
}

// 阻塞取出任务
func (p *Client) BPop(topic string, timeout int) (*Message, error) {
	values, err := redis.Strings(p.Conn.Do("BRPOP", PREFIX_READY_QUEUE+topic, timeout))
	if err != nil {
		return nil, err
	}
	id := values[1]
	result, err := redis.StringMap(p.Conn.Do("HGETALL", PREFIX_JOB_BUCKET+id))
	if err != nil {
		return nil, err
	}
	if result["topic"] == "" || result["body"] == "" {
		return nil, errors.New("Job bucket has expired or is incomplete")
	}
	p.Conn.Do("DEL", PREFIX_JOB_BUCKET+id)
	msg := &Message{
		ID:    id,
		Topic: result["topic"],
		Body:  result["body"],
	}
	return msg, nil
}

// 移除任务
func (p *Client) Remove(id string) (bool, error) {
	// 执行事务
	p.Conn.Send("MULTI")
	p.Conn.Send("ZREM", KEY_JOB_POOL, id)
	p.Conn.Send("DEL", PREFIX_JOB_BUCKET+id)
	values, err := redis.Values(p.Conn.Do("EXEC"))
	if err != nil {
		return false, err
	}
	// 事务结果处理
	v := values[0].(int64)
	v1 := values[1].(int64)
	if v == 0 || v1 == 0 {
		return false, nil
	}
	// 返回
	return true, nil
}
