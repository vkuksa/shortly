package http

import "strconv"

type Config struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

func (c *Config) BuildAddr() string {
	addr := ":" + strconv.Itoa(int(c.Port))

	if c.Host == "" {
		return addr
	}

	return c.Host + addr
}
