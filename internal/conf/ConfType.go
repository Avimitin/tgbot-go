package conf

import "fmt"

type Config struct {
	groups     []int64
	botToken   string
	db         DBSecret
	keywords   []*KeywordsReplyType
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
	var lo, hi int = 0, len(c.groups)
	for mid := lo; mid >= lo && mid <= hi; mid = hi - (hi-lo)/2 {
		if id > c.groups[mid] {
			lo = mid + 1
			continue
		}
		if id < c.groups[mid] {
			hi = mid - 1
			continue
		}
		return true
	}
	return false
}

type Context struct {
	done     chan bool
	errorMSG []string
}

type KeywordsReplyType map[string][]string

// DBSecret store database DSN information
type DBSecret struct {
	user     string
	pwd      string
	host     string
	database string
	port     string
}

// MySqlURL return formatted mysql tcp link
func (db *DBSecret) MySqlURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db.user, db.pwd, db.host, db.port, db.database)
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
