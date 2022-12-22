// package fsdd ...
package fsdd

// Config is the general application layer high-level interface struct
type Config struct {
	Opt                  string // clear text of activated commanddline options
	Path                 string // target path
	HardLink             bool   // activate in-place hardlinks for duplicated data
	SymLink              bool   // activate in-place symlinks for duplicated data
	ReplaceSymlinks      bool   // replace all valid symlinks via hardlinks
	RemoveBrokenSymlinks bool   // delete all broken symlinks
	CleanMetadata        bool   // set mtime/atime to unix-zero [1970-01-01 00:00]
	FastHash             bool   // use xxh3 hash instead of crytographic secure sha512/256 hash
	Verbose              bool   // verbose report about fs state
	Debug                bool
}

// GetDefaultConfig ...
func GetDefaultConfig() *Config {
	return &Config{
		FastHash: true,
	}
}

// Start ...
func (config *Config) Start() {
	config.run()
}
