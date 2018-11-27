package centtest

import (
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

type IDGenerator interface {
	ID() string
}

func NewIDGenerator(source string) IDGenerator {
	switch {
	case strings.HasPrefix(source, "redis://"):
		r, err := redis.DialURL(source)
		if err != nil {
			panic(err)
		}
		return &redisGenerator{r}
	default:
		return &uuidGenerator{}
	}
}

type uuidGenerator struct{}

func (*uuidGenerator) ID() string {
	return uuid.New().String()
}

type redisGenerator struct {
	r redis.Conn
}

func (g *redisGenerator) ID() string {
	n, err := redis.Int(g.r.Do("INCR", "centtest:redisGenerator:ID"))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", n)
}
