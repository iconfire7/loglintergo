package config

type Config struct {
	Rules Rules `mapstructure:"rules"`
}

type Rules struct {
	Lowercase      bool `mapstructure:"lowercase"`
	English        bool `mapstructure:"english"`
	EmojiOrSpesial bool `mapstructure:"emoji_or_spesial"`
	Sensitive      bool `mapstructure:"sensitive"`
}

func Default() Config {
	return Config{
		Rules: Rules{
			Lowercase:      true,
			English:        true,
			EmojiOrSpesial: true,
			Sensitive:      true,
		},
	}
}
