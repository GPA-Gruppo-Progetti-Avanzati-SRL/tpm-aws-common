package awss3lks

type BlobInfo struct {
	Name      string `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Container string `mapstructure:"bucket,omitempty" yaml:"bucket,omitempty" json:"bucket,omitempty"`
	Version   string `mapstructure:"version,omitempty" yaml:"version,omitempty" json:"version,omitempty"`
	Public    bool   `mapstructure:"public,omitempty" yaml:"public,omitempty" json:"public,omitempty"`
}
