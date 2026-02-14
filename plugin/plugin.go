package plugin

import (
	"fmt"
	"regexp"

	"github.com/golangci/plugin-module-register/register"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/tools/go/analysis"

	"github.com/iconfire7/loglintergo/internal/analyzer/loglinter"
	"github.com/iconfire7/loglintergo/internal/config"
)

func init() {
	register.Plugin("loglintergo", New)
}

type Plugin struct {
	cfg       config.Config
	sensitive []*regexp.Regexp
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{loglinter.New(p.cfg, p.sensitive)}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func New(settings any) (register.LinterPlugin, error) {
	cfg := config.Default()

	if settings == nil {
		return &Plugin{cfg: cfg}, nil
	}

	m, ok := settings.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected settings type %T", settings)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName:          "mapstructure",
		Result:           &cfg,
		WeaklyTypedInput: true,
	})
	if err != nil {
		return nil, err
	}
	if err := decoder.Decode(m); err != nil {
		return nil, fmt.Errorf("decode settings: %w", err)
	}

	var reg []*regexp.Regexp
	if cfg.Rules.Sensitive && len(cfg.SensitivePatterns) > 0 {
		reg = make([]*regexp.Regexp, 0, len(cfg.SensitivePatterns))
		for _, p := range cfg.SensitivePatterns {
			re, err := regexp.Compile(p)
			if err != nil {
				return nil, fmt.Errorf("invalid sensitive pattern %q: %w", p, err)
			}
			reg = append(reg, re)
		}
	}

	return &Plugin{cfg: cfg, sensitive: reg}, nil
}
