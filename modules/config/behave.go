package config

// GetOwner return bot instance owner id
func GetOwner() int {
	return cfg.Bot.Owner
}

// GetBotToken return bot instance token
func GetBotToken() string {
	return cfg.Bot.Token
}

// GetDatabaseDSN return a mysql format dsn
func GetDatabaseDSN() string {
	return cfg.Database.EncodeMySQLDSN()
}

// GetDatabaseLogLevel return log level setting for orm
func GetDatabaseLogLevel() string {
	return cfg.Database.LogLevel
}

// GetBotLogLevel return log level setting for bot instance
func GetBotLogLevel() string {
	return cfg.Bot.LogLevel
}
