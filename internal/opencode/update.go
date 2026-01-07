package opencode

import (
	"encoding/json"
	"fmt"
	"os"
)

// UpdatePluginSpec updates the declared plugin spec in the config file.
func UpdatePluginSpec(path string, pluginName string, newSpec string) error {
	if path == "" {
		return fmt.Errorf("config path is required")
	}
	if pluginName == "" {
		return fmt.Errorf("plugin name is required")
	}
	if newSpec == "" {
		return fmt.Errorf("new spec is required")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	var raw map[string]any
	if err := json.Unmarshal(stripJSONC(data), &raw); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}

	updated := false
	if list, ok, err := updateList(raw, "plugin", pluginName, newSpec); err != nil {
		return err
	} else if ok {
		raw["plugin"] = list
		updated = true
	}

	if list, ok, err := updateList(raw, "plugins", pluginName, newSpec); err != nil {
		return err
	} else if ok {
		raw["plugins"] = list
		updated = true
	}

	if !updated {
		return ErrPluginNotFound
	}

	out, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	out = append(out, '\n')
	if err := os.WriteFile(path, out, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func updateList(raw map[string]any, key string, pluginName string, newSpec string) ([]string, bool, error) {
	value, ok := raw[key]
	if !ok {
		return nil, false, nil
	}
	list, err := coerceStringSlice(value)
	if err != nil {
		return nil, false, fmt.Errorf("parse %s list: %w", key, err)
	}

	updated := false
	for i, spec := range list {
		name, _ := parseSpec(spec)
		if name == pluginName {
			list[i] = newSpec
			updated = true
		}
	}
	return list, updated, nil
}

func coerceStringSlice(value any) ([]string, error) {
	switch typed := value.(type) {
	case []string:
		return typed, nil
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			str, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("non-string value")
			}
			out = append(out, str)
		}
		return out, nil
	default:
		return nil, fmt.Errorf("unsupported list type")
	}
}
