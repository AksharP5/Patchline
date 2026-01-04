package cache

type Entry struct {
	Name    string
	Version string
	Path    string
}

func Detect(cacheDir string) ([]Entry, error) {
	return nil, ErrNotImplemented
}
