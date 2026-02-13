package loglintergo_test

import (
	"fmt"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/iconfire7/loglintergo/internal/analyzer"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	fmt.Println("testdata dir:" + testdata)
	analysistest.Run(t, testdata, analyzer.Analyzer, "server")
}
