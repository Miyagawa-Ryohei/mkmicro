package types

type EndPoint struct {
	Region string `toml:"region" mapstructure:"region"`
	URL    string `toml:"url" mapstructure:"url"`
}

type Credential struct {
	AccessKey       string `mapstructure:"access_key,omitempty"`
	AccessKeySecret string `mapstructure:"access_key_secret,omitempty"`
}

type Profile struct {
	Name       string `mapstructure:"name,omitempty"`
	Role       string `mapstructure:"role,omitempty"`
	AssumeRole string `mapstructure:"assume_role,omitempty"`
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
