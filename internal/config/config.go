package config

type HTTP struct {
	Addr string
}

type LastFM struct {
	APIKey string
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
	HTTP   HTTP
	LastFM  LastFM
	KKBOX   KKBOX
	Spotify Spotify
}