package httpserver

//goland:noinspection GoUnusedGlobalVariable
var DefaultConfig = Config{
	URL:  "/interactions/callback",
	Port: ":80",
}

type Config struct {
	URL       string
	Port      string
	PublicKey string
	CertFile  string
	KeyFile   string
}

type ConfigOpt func(config *Config)

func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithURL(url string) ConfigOpt {
	return func(config *Config) {
		config.URL = url
	}
}

func WithPort(port string) ConfigOpt {
	return func(config *Config) {
		config.Port = port
	}
}

func WithPublicKey(publicKey string) ConfigOpt {
	return func(config *Config) {
		config.PublicKey = publicKey
	}
}

func WithTLS(certFile string, keyFile string) ConfigOpt {
	return func(config *Config) {
		config.CertFile = certFile
		config.KeyFile = keyFile
	}
}
