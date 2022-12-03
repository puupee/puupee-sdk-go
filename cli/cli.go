package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/puupee/puupee-api-go"
	"github.com/spf13/viper"
)

type Session struct {
	AccessToken string `json:"access_token" yaml:"access_token" mapstructure:"access_token"`
	ExpiresIn   int64  `json:"expires_in" yaml:"expires_in" mapstructure:"expires_in"`
	CreatedAt   int64  `json:"created_at" yaml:"created_at" mapstructure:"created_at"`
	TokenType   string `json:"token_type" yaml:"token_type" mapstructure:"token_type"`
	Scope       string `json:"scope" yaml:"scope" mapstructure:"scope"`
	PhoneNumber string `json:"phone_number" yaml:"phone_number" mapstructure:"phone_number"`
}

type Config struct {
	Host    string   `json:"host" yaml:"host" mapstructure:"host"`
	Session *Session `json:"session" yaml:"session" mapstructure:"session"`
}

func NewConfig() *Config {
	return &Config{
		Host:    "api.puupee.code2code.cn",
		Session: &Session{},
	}
}

func (session *Session) Valid() bool {
	if session.AccessToken == "" {
		return false
	}
	if session.ExpiresIn <= 0 {
		return false
	}
	if session.CreatedAt <= 0 {
		return false
	}
	return session.CreatedAt+session.ExpiresIn > time.Now().Unix()
}

type puupeeCli struct {
	api     *puupee.APIClient
	session *Session
	config  *Config

	AppOp     *AppOp
	ReleaseOp *ReleaseOp
}

func NewpuupeeCli() *puupeeCli {
	cliCfg := NewConfig()
	err := viper.Unmarshal(&cliCfg)

	puupeeCfg := puupee.NewConfiguration()
	puupeeCfg.Scheme = "https"
	puupeeCfg.Host = cliCfg.Host
	api := puupee.NewAPIClient(puupeeCfg)

	if err != nil {
		panic(err)
	}
	if cliCfg.Session == nil {
		cliCfg.Session = &Session{}
	}
	if cliCfg.Session.Valid() {
		api.GetConfig().AddDefaultHeader("Authorization", cliCfg.Session.TokenType+" "+cliCfg.Session.AccessToken)
	}
	return &puupeeCli{
		api:       api,
		session:   cliCfg.Session,
		config:    cliCfg,
		AppOp:     NewAppOp(api),
		ReleaseOp: NewReleaseOp(api),
	}
}

func (cli *puupeeCli) Login() error {
	if cli.session.Valid() {
		return fmt.Errorf("已经登录， 无需重复登录")
	}
	var phoneNumber string
	if err := survey.AskOne(&survey.Input{Message: "请输入手机号码+86:"}, &phoneNumber); err != nil {
		return err
	}
	if !strings.HasPrefix("+86", phoneNumber) {
		phoneNumber = "+86" + phoneNumber
	}
	_, err := cli.api.VerificationApi.ApiAppVerificationSendCodePost(context.Background()).
		Body(puupee.SendVerificationCodeDto{
			Account: &phoneNumber,
		}).Execute()
	if err != nil {
		return err
	}
	var smsCode string
	if err := survey.AskOne(&survey.Input{Message: "请输入短信验证码:"}, &smsCode); err != nil {
		return err
	}
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	deviceToken := id
	v := url.Values{}
	v.Set("grant_type", "sms")
	v.Set("scope", "puupee")
	v.Set("client_id", "puupee_Sms_GrantType")
	v.Set("client_secret", "1q2w3e*")
	v.Set("phone_number", phoneNumber)
	v.Set("sms_code", smsCode)
	v.Set("device_token", deviceToken)
	v.Set("device_name", "puupee-cli")
	v.Set("device_platform_type", "other")
	v.Set("device_brand", "puupee-cli")
	v.Set("device_system_version", "1.0.0")

	rsp, err := cli.api.GetConfig().HTTPClient.PostForm(cli.api.GetConfig().Scheme+"://"+cli.api.GetConfig().Host+"/connect/token", v)
	if err != nil {
		return err
	}
	if rsp.StatusCode > 300 {
		return fmt.Errorf("登录失败，请检查手机号码和短信验证码")
	}
	bts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(bts))
	session := &Session{}
	err = json.Unmarshal(bts, &session)
	if err != nil {
		return err
	}
	cli.session = session
	session.PhoneNumber = phoneNumber
	session.CreatedAt = time.Now().Unix()
	viper.Set("session", session)
	return viper.WriteConfig()
}

func (cli *puupeeCli) Logout() error {
	viper.Set("session", nil)
	return viper.WriteConfig()
}
