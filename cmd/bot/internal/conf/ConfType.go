package conf

type dbCfg struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

type Config struct {
	LOADED   bool `yaml:"LOADED"`
	Groups   []int64
	BotToken string `yaml:"bot_token"`
	DBCfg    dbCfg  `yaml:"db_cfg"`
}
