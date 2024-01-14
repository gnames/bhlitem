package config

type Config struct {
	Path string
}

type Option func(*Config)

func OptPath(s string) Option {
	return func(cfg *Config) {
		cfg.Path = s
	}
}

func New(opts ...Option) Config {
	res := Config{
		Path: "/opt/bhl",
	}

	for _, opt := range opts {
		opt(&res)
	}
	return res
}
