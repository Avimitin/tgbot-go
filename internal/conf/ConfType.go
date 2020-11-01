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
	Keywords []*KeywordsReplyType
}

type KeywordsReplyType map[string][]string

type DBSecret struct {
	User     string
	Pwd      string
	Host     string
	Database string
	Port     string
}

// BotData store map of kw and replies and auth groups
type BotData struct {
	KAR    KeywordsReplyType
	Groups []int64
}

// BDInit Init all the data in BotData type
func (b *BotData) BDInit() {
	b.KAR = make(KeywordsReplyType)
	b.Groups = make([]int64, 5)
}
