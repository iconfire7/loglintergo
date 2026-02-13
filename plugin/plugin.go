package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/iconfire7/loglintergo/internal/analyzer"
)

type Plugin struct{}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}

func New(_ any) (register.LinterPlugin, error) {
	return &Plugin{}, nil
}

func init() {
	register.Plugin("loglintergo", New)
}
