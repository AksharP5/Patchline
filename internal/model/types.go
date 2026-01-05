package model

type Status string

const (
	StatusOK        Status = "ok"
	StatusMissing   Status = "missing"
	StatusMismatch  Status = "mismatch"
	StatusUnmanaged Status = "unmanaged"
	StatusOutdated  Status = "outdated"
	StatusUnknown   Status = "unknown"
)

type Plugin struct {
	Name           string
	DeclaredSpec   string
	Installed      string
	Status         Status
	Source         string
	ConfigPath     string
	CachePath      string
	LocalDirectory string
}
