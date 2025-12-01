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

type LastFM struct {
	APIKey string
}

type YTMusic struct {
	SidecarURL string
}

type Config struct {
	HTTP    HTTP
	KKBOX   KKBOX
	Spotify Spotify
	LastFM  LastFM
	YTMusic YTMusic
}
