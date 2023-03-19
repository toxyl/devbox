package config

type Commands struct {
	Start string `mapstructure:"start"`
	Stop  string `mapstructure:"stop"`
}
