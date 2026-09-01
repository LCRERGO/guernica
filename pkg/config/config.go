package config

type Config struct {
	Length           int
	PasswordAlphabet string
	Count            int
	Loop             bool
	Clipboard        bool
	Words            int
	Separator        string
	WordlistFile     string
	NoEntropy        bool
	RequireAll       bool
	NoAmbiguous      bool
	AmbiguousExclude string
	CustomAlphabet   string
}

var defConfig *Config

func GetConfig() *Config {
	if defConfig == nil {
		defConfig = &Config{}
	}

	return defConfig
}
