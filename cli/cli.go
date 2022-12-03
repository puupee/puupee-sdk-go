package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/puupee/puupee-api-go"
	"github.com/spf13/viper"
)

type Session struct {
	AccessToken  string `json:"access_token" yaml:"access_token" mapstructure:"access_token"`
	IdToken      string `json:"id_token" yaml:"id_token" mapstructure:"id_token"`
	RefreshToken string `json:"refresh_token" yaml:"refresh_token" mapstructure:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in" yaml:"expires_in" mapstructure:"expires_in"`
	CreatedAt    int64  `json:"created_at" yaml:"created_at" mapstructure:"created_at"`
	TokenType    string `json:"token_type" yaml:"token_type" mapstructure:"token_type"`
	Scope        string `json:"scope" yaml:"scope" mapstructure:"scope"`
	PhoneNumber  string `json:"phone_number" yaml:"phone_number" mapstructure:"phone_number"`
	Username     string `json:"username" yaml:"username" mapstructure:"username"`
}

type Config struct {
	Host    string   `json:"host" yaml:"host" mapstructure:"host"`
	Session *Session `json:"session" yaml:"session" mapstructure:"session"`
}

func NewConfig() *Config {
	return &Config{
		Host:    "api.puupee.com",
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
		authorization := fmt.Sprintf("%s %s", cliCfg.Session.TokenType, cliCfg.Session.AccessToken)
		api.GetConfig().AddDefaultHeader("Authorization", authorization)
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
	loginMethodPrompt := &survey.Select{
		Message: "请选择登录方式:",
		Options: []string{"账号密码", "短信验证码"},
	}
	var loginMethod string
	var phoneNumber string
	var smsCode string
	var username string
	var password string
	err := survey.AskOne(loginMethodPrompt, &loginMethod, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}
	if loginMethod == "账号密码" {
		if err := survey.AskOne(&survey.Input{Message: "请输入用户名:"}, &username); err != nil {
			return err
		}
		if err := survey.AskOne(&survey.Password{Message: "请输入密码:"}, &password); err != nil {
			return err
		}
	}
	codeSender := "SMS"
	loginPurpose := "Login"
	if loginMethod == "短信验证码" {
		if err := survey.AskOne(&survey.Input{Message: "请输入手机号码+86:"}, &phoneNumber); err != nil {
			return err
		}
		if !strings.HasPrefix("+86", phoneNumber) {
			phoneNumber = "+86" + phoneNumber
		}
		_, err := cli.api.VerificationApi.ApiAppVerificationSendCodePost(context.Background()).
			Body(puupee.SendVerificationCodeDto{
				CodeSender: &codeSender,
				Account:    &phoneNumber,
				Purpose:    &loginPurpose,
			}).Execute()
		if err != nil {
			return err
		}
		if err := survey.AskOne(&survey.Input{Message: "请输入短信验证码:"}, &smsCode); err != nil {
			return err
		}
	}

	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	deviceToken := id
	v := url.Values{}
	puupeeClientId := os.Getenv("PUUPEE_CLIENT_ID")
	puupeeClientSecret := os.Getenv("PUUPEE_CLIENT_SECRET")

	v.Set("grant_type", "phone_sms_verify")
	v.Set("scope", "Puupees openid offline_access address email phone profile roles")
	v.Set("client_id", puupeeClientId)
	v.Set("client_secret", puupeeClientSecret)
	v.Set("phone_number", phoneNumber)
	v.Set("sms_code", smsCode)
	v.Set("username", username)
	v.Set("password", password)
	v.Set("device_token", deviceToken)
	v.Set("device_name", "puupee-cli")
	v.Set("device_platform_type", "command-line")
	v.Set("device_brand", "puupee-cli")
	v.Set("device_system_version", "1.0.0")

	loginUrl := fmt.Sprintf("%s://%s/connect/token", cli.api.GetConfig().Scheme, cli.api.GetConfig().Host)
	rsp, err := cli.api.GetConfig().HTTPClient.PostForm(loginUrl, v)
	if err != nil {
		return err
	}
	bts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(bts))
	if rsp.StatusCode > 300 {
		return fmt.Errorf("登录失败，请检查手机号码和短信验证码")
	}
	session := &Session{}
	err = json.Unmarshal(bts, &session)
	if err != nil {
		return err
	}
	cli.session = session
	session.Username = username
	session.PhoneNumber = phoneNumber
	session.CreatedAt = time.Now().Unix()
	viper.Set("session", session)
	return viper.WriteConfig()
}

func (cli *puupeeCli) Logout() error {
	viper.Set("session", nil)
	return viper.WriteConfig()
}
