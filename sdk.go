package puupeesdk

import (
	"github.com/puupee/puupee-api-go"
	"github.com/spf13/viper"
)

type Config struct {
	Env     string            `json:"env,omitempty" yaml:"env" mapstructure:"env"`
	Host    string            `json:"host,omitempty" yaml:"host" mapstructure:"host"`
	ApiKey  string            `json:"apiKey,omitempty" yaml:"apiKey" mapstructure:"apiKey"`
	ApiKeys map[string]string `json:"apiKeys,omitempty" yaml:"apiKeys" mapstructure:"apiKeys"`
}

func (c *Config) GetApiKey() string {
	if c.ApiKey != "" {
		return c.ApiKey
	}
	return c.ApiKeys[c.Env]
}

func (c *Config) SetApiKey(name string, value string) {
	c.ApiKey = value
	c.ApiKeys[name] = value
}

func NewConfig() *Config {
	return &Config{
		Env:     "prod",
		Host:    "api.puupee.com",
		ApiKeys: map[string]string{},
	}
}

type PuupeeSdk struct {
	api    *puupee.APIClient
	config *Config

	App     *AppOp
	Release *ReleaseOp
	ApiKey  *ApiKeyOp
}

func NewSdk() *PuupeeSdk {
	cliCfg := NewConfig()
	err := viper.Unmarshal(&cliCfg)

	puupeeCfg := puupee.NewConfiguration()
	puupeeCfg.Scheme = "https"
	puupeeCfg.Host = cliCfg.Host
	puupeeCfg.DefaultHeader["X-Requested-With"] = "XMLHttpRequest"
	puupeeCfg.DefaultHeader["api-key"] = cliCfg.GetApiKey()

	api := puupee.NewAPIClient(puupeeCfg)

	if err != nil {
		panic(err)
	}

	cli := &PuupeeSdk{
		api:     api,
		config:  cliCfg,
		App:     NewAppOp(api),
		Release: NewReleaseOp(api),
		ApiKey:  NewApiKeyOp(api),
	}
	return cli
}
