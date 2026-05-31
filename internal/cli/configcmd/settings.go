package configcmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/shhac/lin/internal/config"
	"github.com/shhac/lin/internal/output"
)

type settingDef struct {
	description string
	get         func(*config.Settings) any
	apply       func(*config.Settings, string) (any, error)
	reset       func(*config.Settings)
}

var settingDefs = map[string]settingDef{
	"truncation.maxLength": {
		description: "Max characters before truncating description/body/content fields (default: 200)",
		apply: func(s *config.Settings, raw string) (any, error) {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 0 {
				return nil, fmt.Errorf("invalid value: %s, must be a non-negative integer", raw)
			}
			if s.Truncation == nil {
				s.Truncation = &config.TruncationSettings{}
			}
			s.Truncation.MaxLength = &n
			return n, nil
		},
		get: func(s *config.Settings) any {
			if s.Truncation == nil || s.Truncation.MaxLength == nil {
				return nil
			}
			return *s.Truncation.MaxLength
		},
		reset: func(s *config.Settings) {
			if s.Truncation != nil {
				s.Truncation.MaxLength = nil
			}
		},
	},
	"pagination.defaultPageSize": {
		description: "Default number of results for list/search commands (default: 50)",
		apply: func(s *config.Settings, raw string) (any, error) {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 1 || n > 250 {
				return nil, fmt.Errorf("invalid value: %s, must be an integer between 1 and 250", raw)
			}
			if s.Pagination == nil {
				s.Pagination = &config.PaginationSettings{}
			}
			s.Pagination.DefaultPageSize = &n
			return n, nil
		},
		get: func(s *config.Settings) any {
			if s.Pagination == nil || s.Pagination.DefaultPageSize == nil {
				return nil
			}
			return *s.Pagination.DefaultPageSize
		},
		reset: func(s *config.Settings) {
			if s.Pagination != nil {
				s.Pagination.DefaultPageSize = nil
			}
		},
	},
	"output.defaultFormat": {
		description: "Default output format when --format is omitted (json, yaml, jsonl)",
		apply: func(s *config.Settings, raw string) (any, error) {
			if _, err := output.ParseFormat(raw); err != nil {
				return nil, err
			}
			value := raw
			if value == "ndjson" {
				value = "jsonl"
			}
			if s.Output == nil {
				s.Output = &config.OutputSettings{}
			}
			s.Output.DefaultFormat = value
			return value, nil
		},
		get: func(s *config.Settings) any {
			if s.Output == nil || s.Output.DefaultFormat == "" {
				return nil
			}
			return s.Output.DefaultFormat
		},
		reset: func(s *config.Settings) {
			if s.Output != nil {
				s.Output.DefaultFormat = ""
			}
		},
	},
	"request.timeoutMS": {
		description: "Default request timeout in milliseconds (0 disables client timeout)",
		apply: func(s *config.Settings, raw string) (any, error) {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 0 {
				return nil, fmt.Errorf("invalid value: %s, must be a non-negative integer", raw)
			}
			if s.Request == nil {
				s.Request = &config.RequestSettings{}
			}
			s.Request.TimeoutMS = &n
			return n, nil
		},
		get: func(s *config.Settings) any {
			if s.Request == nil || s.Request.TimeoutMS == nil {
				return nil
			}
			return *s.Request.TimeoutMS
		},
		reset: func(s *config.Settings) {
			if s.Request != nil {
				s.Request.TimeoutMS = nil
			}
		},
	},
}

var validKeys = func() []string {
	keys := make([]string, 0, len(settingDefs))
	for k := range settingDefs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}()

func validKeysStr() string {
	return strings.Join(validKeys, ", ")
}
