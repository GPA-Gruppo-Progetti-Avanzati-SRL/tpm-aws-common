package awss3lks

type BucketCfg struct {
	Tag     string   `yaml:"tag,omitempty" mapstructure:"tag,omitempty" json:"tag,omitempty"`
	Buckets []string `yaml:"buckets,omitempty" mapstructure:"buckets,omitempty" json:"buckets,omitempty"`
	Path    string   `yaml:"path,omitempty" mapstructure:"path,omitempty" json:"path,omitempty"`
}

type Config struct {
	Name               string      `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Endpoint           string      `mapstructure:"endpoint,omitempty" yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	AccessKey          string      `mapstructure:"access-key,omitempty" yaml:"access-key,omitempty" json:"access-key,omitempty"`
	SecretKey          string      `mapstructure:"secret-key,omitempty"  yaml:"secret-key,omitempty" json:"secret-key,omitempty"`
	Region             string      `mapstructure:"region,omitempty"  yaml:"region,omitempty" json:"region,omitempty"`
	PublicEndpoint     string      `mapstructure:"public-url,omitempty"  yaml:"public-url,omitempty" json:"public-url,omitempty"`
	BucketConfig       []BucketCfg `mapstructure:"buckets,omitempty"  yaml:"buckets,omitempty" json:"buckets,omitempty"`
	UseSharedAWSConfig bool        `mapstructure:"with-shared-cfg,omitempty"  yaml:"with-shared-cfg,omitempty" json:"with-shared-cfg,omitempty"`
	UseGoogleConfig    bool        `mapstructure:"use-google-config,omitempty"  yaml:"use-google-config,omitempty" json:"use-google-config,omitempty"`
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
