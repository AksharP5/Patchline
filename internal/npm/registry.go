package npm

type PackageInfo struct {
	Name    string
	Latest  string
	Versions []string
}

func FetchPackageInfo(name string) (PackageInfo, error) {
	return PackageInfo{}, ErrNotImplemented
}
