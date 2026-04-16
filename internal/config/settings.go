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
	cfg := Read()
	if cfg.Settings == nil {
		cfg.Settings = &Settings{}
	}
	if partial.Truncation != nil {
		cfg.Settings.Truncation = partial.Truncation
	}
	if partial.Pagination != nil {
		cfg.Settings.Pagination = partial.Pagination
	}
	return Write(cfg)
}

// ResetSettings removes all settings.
func ResetSettings() error {
	cfg := Read()
	cfg.Settings = nil
	return Write(cfg)
}
