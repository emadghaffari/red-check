package authenticator

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

const defaultDB = 0
const defaultUsername = ""
const defaultHost = ""
const defaultPassword = ""

type Config struct {
	DB       int
	Host     string
	Username string
	Password string
}

type Authenticator struct {
	DB       int
	Username string
	Password string
	Host     string
	err      error
	client   *redis.Client
}

type Result struct {
	value string
	exist bool
	err   error
}

func New() Authenticator {
	return Authenticator{
		DB:       defaultDB,
		Username: defaultUsername,
		Password: defaultPassword,
		Host:     defaultHost,
	}
}

func (a *Authenticator) SetConfig(config Config) *Authenticator {
	a.DB = config.DB
	a.Host = config.Host
	a.Username = config.Username
	a.Password = config.Password
	return a
}

func (a *Authenticator) Connect() *Authenticator {
	a.client = redis.NewClient(&redis.Options{
		DB:       a.DB,
		Addr:     a.Host,
		Username: a.Username,
		Password: a.Password,
	})

	if err := a.client.Ping(context.Background()).Err(); err != nil {
		a.err = err
	}

	return a
}

func (a *Authenticator) Get(ctx context.Context, key string) Result {
	var resp Result
	p, err := a.client.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return Result{err: err}
	}

	resp.exist = true
	err = json.Unmarshal([]byte(p), &resp.value)
	if err != nil {
		return Result{}
	}

	return resp
}

func (a *Authenticator) IsExist(ctx context.Context, key string) (bool, error) {
	i, err := a.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return (i == 1), nil

}

func (a *Authenticator) Error() error {
	return a.err
}

func (r Result) Exist() bool {
	return r.exist
}

func (r Result) Error() error {
	return r.err
}

func (r Result) Value() string {
	return r.value
}
