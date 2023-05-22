package http

type Config struct {
	Scheme string `toml:"scheme"`
	Host   string `toml:"host"`
	Port   string `toml:"port"`

	Prometheus struct {
		Enabled bool   `toml:"enabled"`
		Port    string `toml:"port"`
	} `toml:"prometheus"`
}
