package config

type Config struct {
	Length           int
	PasswordAlphabet string
}

var defConfig *Config

func GetConfig() *Config {
	if defConfig == nil {
		defConfig = &Config{}
	}

	return defConfig
}
