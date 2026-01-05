package snapshot

import "time"

type Entry struct {
	Timestamp        time.Time `json:"timestamp"`
	PluginName       string    `json:"pluginName"`
	PreviousSpec     string    `json:"previousSpec"`
	PreviousInstalled string   `json:"previousInstalled"`
	Source           string    `json:"source"`
	Reason           string    `json:"reason"`
}

type Store struct {
	Directory string
}

func (s Store) Save(entry Entry) error {
	return ErrNotImplemented
}

func (s Store) Latest(pluginName string) (Entry, error) {
	return Entry{}, ErrNotImplemented
}
