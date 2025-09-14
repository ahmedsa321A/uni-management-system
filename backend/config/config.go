package config

type Config struct {
	DBPath string
}

func Load() *Config {
	cfg := &Config{
		DBPath: "C:\\projects\\uni_managment\\uni-management-system\\database\\uniDB.db",
	}
	return cfg
}
