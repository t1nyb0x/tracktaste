package config

type HTTP struct {
	Addr string
}

type KKBOX struct {
	APIKey string
	Secret string
}

type Spotify struct {
	APIKey string
	Secret string
}

type Config struct {
	HTTP    HTTP
	KKBOX   KKBOX
	Spotify Spotify
}
