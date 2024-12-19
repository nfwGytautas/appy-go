package appy_config

// SSL settings
type SSLSettings struct {
	CertFile string
	KeyFile  string
}

// HTTP server config
type HttpConfig struct {
	Address string
	SSL     *SSLSettings

	ErrorMapper HttpErrorMapper
}

// Configure the appy app
type AppyConfig struct {
	Http *HttpConfig
}
