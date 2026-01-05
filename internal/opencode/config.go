package opencode

type Source string

const (
	SourceProject Source = "project"
	SourceGlobal  Source = "global"
	SourceLocal   Source = "local"
)

type Config struct {
	Plugins    []string `json:"plugin"`
	PluginsAlt []string `json:"plugins"`
}

type PluginSpec struct {
	Name         string
	DeclaredSpec string
	Pinned       string
	Source       Source
	ConfigPath   string
	LocalPath    string
}

type DiscoveryResult struct {
	Plugins []PluginSpec
}
