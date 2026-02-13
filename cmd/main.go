package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/iconfire7/loglintergo/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
