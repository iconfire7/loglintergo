package rules

import (
	"regexp"
	"strconv"
)

type SensitivePattern struct {
	ID string
	Re *regexp.Regexp
}

func CompileSensitive(patterns []string) ([]SensitivePattern, error) {
	out := make([]SensitivePattern, 0, len(patterns))
	for i, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		out = append(out, SensitivePattern{
			ID: "S" + strconv.Itoa(i+1),
			Re: re,
		})
	}
	return out, nil
}
