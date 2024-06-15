package redis

import (
	redigo "github.com/gomodule/redigo/redis"
	"time"
)

// NetworkTCP for tcp.
const NetworkTCP = "tcp"

// Redis defines a redis agent.
type Redis struct {
	pool *redigo.Pool
}

// RedisOptions are options for the redis agent.
type RedisOptions struct {
	MaxIdle   int
	MaxActive int
	Timeout   int
	Wait      bool
}

// RedisConfig config for redis agent.
type RedisConfig struct {
	Address  string
	Password string
	AppName  string
	Network  string
	Options  []RedisOptions
}

// dialFunc is func used to connect to the redis server.
type dialFunc func(network, addr string, opts ...redigo.DialOption) (redigo.Conn, error)

// NewRedis creates new redis agent object.
func NewRedis(cfg RedisConfig) *Redis {
	return newRedisWithDialer(cfg, redigo.Dial)
}

func newRedisWithDialer(cfg RedisConfig, dialer dialFunc) *Redis {
	// Default setting
	opt := RedisOptions{
		MaxIdle:   100,
		MaxActive: 100,
		Timeout:   100,
		Wait:      true,
	}

	if len(cfg.Options) > 0 {
		opt = cfg.Options[0]
	}

	if cfg.Network == "" {
		cfg.Network = NetworkTCP
	}

	return &Redis{
		pool: &redigo.Pool{
			MaxIdle:     opt.MaxIdle,
			MaxActive:   opt.MaxActive,
			IdleTimeout: time.Duration(opt.Timeout) * time.Second,
			Dial: func() (redigo.Conn, error) {
				return dialer(cfg.Network, cfg.Address, redigo.DialPassword(cfg.Password))
			},
		},
	}
}

func (r *Redis) Ping() error {
	conn := r.pool.Get()
	defer conn.Close()

	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Set(key string, value string, ttl int64, field ...interface{}) (interface{}, error) {
	conn := r.pool.Get()
	defer conn.Close()

	hash := len(field) > 0

	defer r.SetExpire(key, ttl)

	if hash {
		return redigo.Int64(conn.Do("HSET", key, field[0], value))
	}
	return redigo.String(conn.Do("SET", key, value))
}

func (r *Redis) Get(key string, field ...interface{}) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	hash := len(field) > 0

	if hash {
		return redigo.String(conn.Do("HGET", key, field[0]))
	}

	return redigo.String(conn.Do("GET", key))
}

func (r *Redis) Del(key string, field ...interface{}) (bool, error) {
	conn := r.pool.Get()
	defer conn.Close()

	hash := len(field) > 0

	if hash {
		return redigo.Bool(conn.Do("HDEL", key, field[0]))
	}
	return redigo.Bool(conn.Do("DEL", key))
}

func (r *Redis) SetExpire(key string, ttl int64) (int64, error) {
	conn := r.pool.Get()
	defer conn.Close()

	if ttl > 0 {
		return redigo.Int64(conn.Do("EXPIRE", key, ttl))
	}
	return 0, nil
}
