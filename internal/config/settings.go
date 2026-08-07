package config

// GetSettings returns the current settings (never nil sub-structs in response).
func GetSettings() *Settings {
	cfg := Read()
	if cfg.Settings == nil {
		return &Settings{}
	}
	return cfg.Settings
}

// UpdateSettings merges partial settings into the config.
func UpdateSettings(partial *Settings) error {
	return update(func(cfg *Config) error {
		if cfg.Settings == nil {
			cfg.Settings = &Settings{}
		}
		if partial.Truncation != nil {
			cfg.Settings.Truncation = partial.Truncation
		}
		if partial.Pagination != nil {
			cfg.Settings.Pagination = partial.Pagination
		}
		if partial.Output != nil {
			cfg.Settings.Output = partial.Output
		}
		if partial.Request != nil {
			cfg.Settings.Request = partial.Request
		}
		return nil
	})
}

// ResetSettings removes all settings.
func ResetSettings() error {
	return update(func(cfg *Config) error {
		cfg.Settings = nil
		return nil
	})
}
