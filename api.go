// {ackage fsdd allows ypu to deduplicate data via hardlinks
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
	FastHash             bool   // use (fast) MAPHASH instead of crytographic secure SHA512/256
	Verbose              bool   // verbose report about fs state
	Debug                bool
}

// DefaultConfig provides as sane default config setup
func DefaultConfig() *Config {
	return &Config{}
}

// Start will perform the requested action from config
func (config *Config) Start() {
	config.run()
}
