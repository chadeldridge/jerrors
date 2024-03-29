package jerrors

var config Config

func init() {
	SetConfig(DefaultConfig())
}

type Config struct {
	LogLevel     bool
	LogTime      bool
	LoggingLevel Level
	LogCaller    bool
	// CallerDepth is how many function calls to step back before getting Caller information.
	// This should be enough to get us back to the initiating function.
	CallerDepth int
	// CallersToShow sets how many calling functions to show.
	CallersToShow int
}

func GetConfig() Config { return config }

func NewConfig() Config {
	return DefaultConfig()
}

func DefaultConfig() Config {
	return Config{
		LogLevel:      true,
		LogTime:       true,
		LoggingLevel:  INFO,
		LogCaller:     false,
		CallerDepth:   2,
		CallersToShow: 2,
	}
}

func SetConfig(newConfig Config) {
	if newConfig.LoggingLevel == 0 {
		newConfig.LoggingLevel = INFO
	}

	config = newConfig
}
