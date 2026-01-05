package snapshot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Entry struct {
	Timestamp         time.Time `json:"timestamp"`
	PluginName        string    `json:"pluginName"`
	PreviousSpec      string    `json:"previousSpec"`
	PreviousInstalled string    `json:"previousInstalled"`
	Source            string    `json:"source"`
	Reason            string    `json:"reason"`
	ConfigPath        string    `json:"configPath"`
}

type Store struct {
	Directory string
}

func (s Store) Save(entry Entry) error {
	if s.Directory == "" {
		return fmt.Errorf("%w: directory is empty", ErrInvalidSnapshotDir)
	}
	if entry.PluginName == "" {
		return fmt.Errorf("plugin name is required")
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	if err := os.MkdirAll(s.Directory, 0o755); err != nil {
		return fmt.Errorf("create snapshot dir: %w", err)
	}

	path := s.entryPath(entry.PluginName)
	entries, err := readEntries(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read snapshot: %w", err)
	}
	entries = append(entries, entry)
	if len(entries) > 1 {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.Before(entries[j].Timestamp)
		})
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}
	data = append(data, '\n')
	if err := writeFileAtomic(path, data); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}
	return nil
}

func (s Store) Latest(pluginName string) (Entry, error) {
	if s.Directory == "" {
		return Entry{}, fmt.Errorf("%w: directory is empty", ErrInvalidSnapshotDir)
	}
	if pluginName == "" {
		return Entry{}, fmt.Errorf("plugin name is required")
	}

	path := s.entryPath(pluginName)
	entries, err := readEntries(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Entry{}, ErrSnapshotNotFound
		}
		return Entry{}, fmt.Errorf("read snapshot: %w", err)
	}
	if len(entries) == 0 {
		return Entry{}, ErrSnapshotNotFound
	}

	latest := entries[0]
	for _, entry := range entries[1:] {
		if entry.Timestamp.After(latest.Timestamp) {
			latest = entry
		}
	}
	return latest, nil
}

func (s Store) entryPath(pluginName string) string {
	name := url.PathEscape(pluginName)
	return filepath.Join(s.Directory, name+".json")
}

func readEntries(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".snapshot-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return nil
}
