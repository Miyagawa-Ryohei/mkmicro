package entity

type QueueConfig struct {
	URL string
}

type SessionConfig struct {
	Endpoint string
	Region string
	RoleArn string
}

type Config struct {
	Queue QueueConfig
	Session SessionConfig
}
