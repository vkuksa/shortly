package metrics

import "strconv"

type Config struct {
	Port int `env:"PORT"`
}

func (c *Config) BuildAddr() string {
	return ":" + strconv.Itoa(int(c.Port))
}
