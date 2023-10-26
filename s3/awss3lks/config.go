package awss3lks

type Config struct {
	Name           string `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Endpoint       string `mapstructure:"endpoint,omitempty" yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	AccessKey      string `mapstructure:"access-key,omitempty" yaml:"access-key,omitempty" json:"access-key,omitempty"`
	SecretKey      string `mapstructure:"secret-key,omitempty"  yaml:"secret-key,omitempty" json:"secret-key,omitempty"`
	Region         string `mapstructure:"region,omitempty"  yaml:"region,omitempty" json:"region,omitempty"`
	PublicEndpoint string `mapstructure:"public-url,omitempty"  yaml:"public-url,omitempty" json:"public-url,omitempty"`
}

type Option func(cfg *Config)

func WithName(k string) Option {
	return func(cfg *Config) {
		cfg.Name = k
	}
}

func WithAccessKey(k string) Option {
	return func(cfg *Config) {
		cfg.AccessKey = k
	}
}

func WithSecretKey(k string) Option {
	return func(cfg *Config) {
		cfg.SecretKey = k
	}
}

func WithEndpoint(k string) Option {
	return func(cfg *Config) {
		cfg.Endpoint = k
	}
}
