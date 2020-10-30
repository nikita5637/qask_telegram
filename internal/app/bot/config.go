package bot

//Config ...
type Config struct {
	QaskAddress string `toml:"qask_address"`
	QaskPort    string `toml:"qask_port"`
	LogLevel    string `toml:"log_level"`
	LogFile     string `toml:"log_file"`
}

//NewConfig ...
func NewConfig() *Config {
	return &Config{}
}
