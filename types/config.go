package types

type EndPoint struct {
	Region string `toml:"region" mapstructure:"region"`
	URL    string `toml:"url" mapstructure:"url"`
}

type Credential struct {
	AccessKey       string
	AccessKeySecret string
}

type Profile struct {
	Name       string
	Role       string
	AssumeRole string
}

type AWSConfig struct {
	Endpoint   *EndPoint   `toml:"endpoint,omitempty" mapstructure:"endpoint,omitempty"`
	Credential *Credential `toml:"credential,omitempty" mapstructure:"queue_credential,omitempty"`
	Profile    *Profile    `toml:"profile,omitempty" mapstructure:"queue_profile,omitempty"`
}

type QueueConfig struct {
	Endpoint   *EndPoint   `mapstructure:"endpoint,omitempty"`
	Credential *Credential `mapstructure:"credential,omitempty"`
	Profile    *Profile    `mapstructure:"profile,omitempty"`
	URL        string      `mapstructure:"url"`
}

func (qc QueueConfig) GetAWSConfig() AWSConfig {
	return AWSConfig{
		Endpoint:   qc.Endpoint,
		Credential: qc.Credential,
		Profile:    qc.Profile,
	}
}

type StorageConfig struct {
	Endpoint   *EndPoint   `mapstructure:"endpoint,omitempty"`
	Credential *Credential `mapstructure:"credential,omitempty"`
	Profile    *Profile    `mapstructure:"profile,omitempty"`
}

func (sc StorageConfig) GetAWSConfig() AWSConfig {
	return AWSConfig{
		Endpoint:   sc.Endpoint,
		Credential: sc.Credential,
		Profile:    sc.Profile,
	}
}

type Config struct {
	Queue   QueueConfig   `mapstructure:"queue"`
	Storage StorageConfig `mapstructure:"storage"`
}
