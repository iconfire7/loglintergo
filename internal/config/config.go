package config

type Config struct {
	Rules             Rules    `mapstructure:"rules"`
	SensitivePatterns []string `mapstructure:"sensitive_patterns"`
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
		SensitivePatterns: []string{
			"(?i)\\\\b(token|secret|api[_-]?key)\\\\b\\\\s*[:=]",
		},
	}
}
