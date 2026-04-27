package configcmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shhac/lin/internal/config"
)

type settingDef struct {
	description string
	parse       func(string) (any, error)
	get         func(*config.Settings) any
	set         func(*config.Settings, int)
	reset       func(*config.Settings)
}

var settingDefs = map[string]settingDef{
	"truncation.maxLength": {
		description: "Max characters before truncating description/body/content fields (default: 200)",
		parse: func(v string) (any, error) {
			n, err := strconv.Atoi(v)
			if err != nil || n < 0 {
				return nil, fmt.Errorf("invalid value: %s, must be a non-negative integer", v)
			}
			return n, nil
		},
		get: func(s *config.Settings) any {
			if s.Truncation == nil || s.Truncation.MaxLength == nil {
				return nil
			}
			return *s.Truncation.MaxLength
		},
		set: func(s *config.Settings, v int) {
			if s.Truncation == nil {
				s.Truncation = &config.TruncationSettings{}
			}
			s.Truncation.MaxLength = &v
		},
		reset: func(s *config.Settings) {
			if s.Truncation != nil {
				s.Truncation.MaxLength = nil
			}
		},
	},
	"pagination.defaultPageSize": {
		description: "Default number of results for list/search commands (default: 50)",
		parse: func(v string) (any, error) {
			n, err := strconv.Atoi(v)
			if err != nil || n < 1 || n > 250 {
				return nil, fmt.Errorf("invalid value: %s, must be an integer between 1 and 250", v)
			}
			return n, nil
		},
		get: func(s *config.Settings) any {
			if s.Pagination == nil || s.Pagination.DefaultPageSize == nil {
				return nil
			}
			return *s.Pagination.DefaultPageSize
		},
		set: func(s *config.Settings, v int) {
			if s.Pagination == nil {
				s.Pagination = &config.PaginationSettings{}
			}
			s.Pagination.DefaultPageSize = &v
		},
		reset: func(s *config.Settings) {
			if s.Pagination != nil {
				s.Pagination.DefaultPageSize = nil
			}
		},
	},
}

var validKeys = func() []string {
	keys := make([]string, 0, len(settingDefs))
	for k := range settingDefs {
		keys = append(keys, k)
	}
	return keys
}()

func validKeysStr() string {
	return strings.Join(validKeys, ", ")
}
