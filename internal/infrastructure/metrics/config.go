package metrics

import "strconv"

type Config struct {
	Port int `toml:"port"`
}

func (c *Config) BuildAddr() string {
	return ":" + strconv.Itoa(int(c.Port))
}
