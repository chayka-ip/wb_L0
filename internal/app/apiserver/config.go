package apiserver

type Config struct {
	BindAddr          string  `toml:"bind_addr"`
	LogLevel          string  `toml:"log_level"`
	DatabaseURL       string  `toml:"database_url"`
	CacheSize         int     `toml:"cache_size"`
	NatsClusterId     string  `toml:"nats_cluster_id"`
	NatsClientId      string  `toml:"nats_client_id"`
	NatsPubliserId    string  `toml:"nats_publisher_id"`
	PublisherWorkTime float64 `toml:"pub_work_time"`
	PublisherSendRate int     `toml:"pub_send_rate"`
	BadDataChance     float32 `toml:"bad_data_chance"`
}

type NatsInfo struct {
	ClusterId string
	ClientId  string
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

func (c *Config) GetNatsInfo() NatsInfo {
	return NatsInfo{
		ClusterId: c.NatsClusterId,
		ClientId:  c.NatsClientId,
	}
}
