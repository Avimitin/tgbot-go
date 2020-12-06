package conf

type dbCfg struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Database string `yaml:"database"`
}

type Config struct {
	LOADED     bool `yaml:"LOADED"`
	Groups     []int64
	BotToken   string `yaml:"bot_token"`
	DBCfg      dbCfg  `yaml:"db_cfg"`
	Keywords   []*KeywordsReplyType
	maxThread  int
	ocpyThread int
	ctx        *Context
}

func (c *Config) RemainThread() bool {
	return c.ocpyThread < c.maxThread
}

// SetOcpyThread expect 0 for decreasing thread and 1 for increasing number.
// This will not change any attribute if given op is unexpected number.
func (c *Config) SetOcpyThread(op int) {
	if op == 1 {
		c.ocpyThread++
		return
	}
	if op == 0 {
		c.ocpyThread--
	}
}

func (c *Config) Context() *Context {
	return c.ctx
}

func (c *Config) SetThread(amounts int) {
	c.maxThread = amounts
}

func (c *Config) InGroups(id int64) bool {
	var lo, hi int = 0, len(c.Groups)
	for mid := lo; mid >= lo && mid <= hi; mid = hi - (hi-lo)/2 {
		if id > c.Groups[mid] {
			lo = mid + 1
			continue
		}
		if id < c.Groups[mid] {
			hi = mid - 1
			continue
		}
		return true
	}
	return false
}

type Context struct {
	Done     chan bool
	ErrorMSG []string
}

func (c *Context) AppendError(err string) {
	c.ErrorMSG = append(c.ErrorMSG, err)
}

func (c *Context) LatestError() string {
	errorMsg := c.ErrorMSG[len(c.ErrorMSG)-1]
	c.ErrorMSG = c.ErrorMSG[:len(c.ErrorMSG)-1]
	return errorMsg
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
