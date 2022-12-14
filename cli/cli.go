package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"runtime"
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
	DeviceToken  string `json:"device_token" yaml:"device_token" mapstructure:"device_token"`
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
	return session.AccessToken != ""
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

	cli := &puupeeCli{
		api:       api,
		session:   cliCfg.Session,
		config:    cliCfg,
		AppOp:     NewAppOp(api),
		ReleaseOp: NewReleaseOp(api),
	}
	if err := cli.RefreshToken(); err != nil {
		cli.session.AccessToken = ""
		fmt.Println("Cleaning invalid access token")
		viper.Set("session", cli.session)
		err = viper.WriteConfig()
		if err != nil {
			panic(err)
		}
	}
	if cliCfg.Session.Valid() {
		authorization := fmt.Sprintf("%s %s", cliCfg.Session.TokenType, cliCfg.Session.AccessToken)
		cli.api.GetConfig().AddDefaultHeader("Authorization", authorization)
	}
	return cli
}

func (cli *puupeeCli) RefreshToken() error {
	deviceToken, err := gonanoid.New()
	if err != nil {
		return err
	}
	v := url.Values{}
	puupeeClientId := os.Getenv("PUUPEE_CLIENT_ID")
	puupeeClientSecret := os.Getenv("PUUPEE_CLIENT_SECRET")
	// hostname, err := os.Hostname()
	// if err != nil {
	// 	return err
	// }
	v.Set("grant_type", "refresh_token")
	v.Set("scope", "Puupees openid offline_access address email phone profile roles")
	v.Set("client_id", puupeeClientId)
	v.Set("client_secret", puupeeClientSecret)
	v.Set("refresh_token", cli.session.RefreshToken)
	// v.Set("device_token", deviceToken)
	// v.Set("device_name", hostname)
	// v.Set("device_platform_type", runtime.GOOS)
	// // TODO: ????????????????????????
	// v.Set("device_brand", "puupee-cli")
	// // TODO: ?????????????????????
	// // https://gist.github.com/flxxyz/ae3ef071dc4ffb0c55daedc7f0740611
	// // https://github.com/matishsiao/goInfo
	// v.Set("device_system_version", "1.0.0")

	loginUrl := fmt.Sprintf("%s://%s/connect/token", cli.api.GetConfig().Scheme, cli.api.GetConfig().Host)
	rsp, err := cli.api.GetConfig().HTTPClient.PostForm(loginUrl, v)
	if err != nil {
		return err
	}
	bts, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	// fmt.Println(string(bts))
	if rsp.StatusCode > 300 {
		return fmt.Errorf("Refresh access_token failed")
	}
	session := &Session{}
	err = json.Unmarshal(bts, &session)
	if err != nil {
		return err
	}
	cli.session.CreatedAt = time.Now().Unix()
	cli.session.DeviceToken = deviceToken
	cli.session.RefreshToken = session.RefreshToken
	cli.session.AccessToken = session.AccessToken
	cli.session.ExpiresIn = session.ExpiresIn
	cli.session.TokenType = session.TokenType
	cli.session.IdToken = session.IdToken
	viper.Set("session", cli.session)
	return viper.WriteConfig()
}

func (cli *puupeeCli) Login() error {
	if cli.session.Valid() {
		return fmt.Errorf("??????????????? ??????????????????")
	}
	loginMethodPrompt := &survey.Select{
		Message: "?????????????????????:",
		Options: []string{"????????????", "???????????????"},
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
	if loginMethod == "????????????" {
		if err := survey.AskOne(&survey.Input{Message: "??????????????????:"}, &username); err != nil {
			return err
		}
		if err := survey.AskOne(&survey.Password{Message: "???????????????:"}, &password); err != nil {
			return err
		}
	}
	codeSender := "SMS"
	loginPurpose := "Login"
	if loginMethod == "???????????????" {
		if err := survey.AskOne(&survey.Input{Message: "?????????????????????+86:"}, &phoneNumber); err != nil {
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
		if err := survey.AskOne(&survey.Input{Message: "????????????????????????:"}, &smsCode); err != nil {
			return err
		}
	}

	deviceToken, err := gonanoid.New()
	if err != nil {
		return err
	}
	v := url.Values{}
	puupeeClientId := os.Getenv("PUUPEE_CLIENT_ID")
	puupeeClientSecret := os.Getenv("PUUPEE_CLIENT_SECRET")
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	v.Set("grant_type", "phone_sms_verify")
	v.Set("scope", "Puupees openid offline_access address email phone profile roles")
	v.Set("client_id", puupeeClientId)
	v.Set("client_secret", puupeeClientSecret)
	v.Set("phone_number", phoneNumber)
	v.Set("sms_code", smsCode)
	v.Set("username", username)
	v.Set("password", password)
	v.Set("device_token", deviceToken)
	v.Set("device_name", hostname)
	v.Set("device_platform_type", runtime.GOOS)
	// TODO: ????????????????????????
	v.Set("device_brand", "puupee-cli")
	// TODO: ?????????????????????
	// https://gist.github.com/flxxyz/ae3ef071dc4ffb0c55daedc7f0740611
	// https://github.com/matishsiao/goInfo
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
	// fmt.Println(string(bts))
	if rsp.StatusCode > 300 {
		return fmt.Errorf("??????????????????????????????????????????????????????")
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
	session.DeviceToken = deviceToken
	viper.Set("session", session)
	return viper.WriteConfig()
}

func (cli *puupeeCli) Logout() error {
	cli.session.AccessToken = ""
	cli.session.RefreshToken = ""
	viper.Set("session", cli.session)
	return viper.WriteConfig()
}
