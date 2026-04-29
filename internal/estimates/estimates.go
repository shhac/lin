package estimates

import (
	"fmt"
	"strings"
)

type Config struct {
	Type      string `json:"type"`
	AllowZero bool   `json:"allowZero"`
	Extended  bool   `json:"extended"`
}

type scale struct {
	base     []int
	extended []int
}

var scales = map[string]scale{
	"exponential": {base: []int{1, 2, 4, 8, 16}, extended: []int{32, 64}},
	"fibonacci":   {base: []int{1, 2, 3, 5, 8, 13}, extended: []int{21, 34}},
	"linear":      {base: []int{1, 2, 3, 4, 5}, extended: []int{6, 7, 8, 9, 10}},
	"tShirt":      {base: []int{1, 2, 3, 4, 5}, extended: []int{6}},
}

var tShirtLabels = map[int]string{
	0: "None", 1: "XS", 2: "S", 3: "M", 4: "L", 5: "XL", 6: "XXL",
}

func BuildConfig(estimationType string, allowZero, extended bool) Config {
	return Config{Type: estimationType, AllowZero: allowZero, Extended: extended}
}

func ValidEstimates(cfg Config) []int {
	s, ok := scales[cfg.Type]
	if !ok {
		return nil
	}
	var values []int
	if cfg.AllowZero {
		values = append(values, 0)
	}
	values = append(values, s.base...)
	if cfg.Extended {
		values = append(values, s.extended...)
	}
	return values
}

func Validate(cfg Config, estimate int) error {
	if cfg.Type == "notUsed" {
		return fmt.Errorf("this team does not use estimates")
	}
	values := ValidEstimates(cfg)
	for _, v := range values {
		if v == estimate {
			return nil
		}
	}
	return fmt.Errorf("invalid estimate %d. Valid values: %s", estimate, FormatScale(cfg.Type, values))
}

func FormatScale(scaleType string, values []int) string {
	parts := make([]string, len(values))
	for i, v := range values {
		if scaleType == "tShirt" {
			label := tShirtLabels[v]
			if label == "" {
				label = fmt.Sprintf("%d", v)
			}
			parts[i] = fmt.Sprintf("%d (%s)", v, label)
		} else {
			parts[i] = fmt.Sprintf("%d", v)
		}
	}
	return strings.Join(parts, " | ")
}
