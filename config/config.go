package config

import (
	"sync"

	"github.com/spf13/viper"
)

type Application struct {
	Name   string `yaml:"name"`
	Server string `yaml:"server"`
}

type Webserver struct {
	Port      uint `yaml:"port"`
	Timeout   uint `yaml:"timeout"`
	RateLimit uint `yaml:"ratelimit"`
}

type Database struct {
	Tipe        string `yaml:"tipe"`
	Driver      string `yaml:"driver"`
	Address     string `yaml:"address"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Port        uint   `yaml:"port"`
	Name        string `yaml:"name"`
	MaxIdle     uint   `yaml:"maxidle"`
	MaxOpen     uint   `yaml:"maxopen"`
	MaxIdleTime uint   `yaml:"maxidletime"`
	MaxLifeTime uint   `yaml:"maxlifetime"`
	Timeout     uint   `yaml:"timeout"`
}

type Jwt struct {
	VerifyKey               string `yaml:"verifykey"`
	Key                     string `yaml:"key"`
	Tokenexpiredtime        uint   `yaml:"tokenexpiredtime"`
	Refreshtokenexpiredtime uint   `yaml:"refreshtokenexpiredtime"`
	FormVerifyKey           string `yaml:"formverifykey"`
	FormKey                 string `yaml:"formkey"`
	FormTokenExperiedTime   uint   `yaml:"formtokenexpiredtime"`
}

type Timezone struct {
	Timezone string `yaml:"timezone"`
}

type Log struct {
	Level  string   `json:"level"`
	Levels []string `json:"Levels"`
}

type IpaymuPayment struct {
	IpaymuVa             string `yaml:"ipaymuva"`
	IpaymuKey            string `yaml:"ipaymukey"`
	IpaymuUrl            string `yaml:"ipaymuurl"`
	IpaymuSnapUrl        string `yaml:"ipaymusnapurl"`
	IpaymuCallbackUrl    string `yaml:"ipaymucallbackurl"`
	IpaymuTranscationUrl string `yaml:"ipaymutranscationurl"`
	IpaymuThankYouPage   string `yaml:"ipaymuthankyoupage"`
	IpaymuCancelUrl      string `yaml:"ipaymucancelurl"`
}

type Email struct {
	FromEmail         string `yaml:"fromemail"`
	FromEmailPassword string `yaml:"fromemailpassword"`
	LinkVerifyEmail   string `yaml:"linkverifyemail"`
}

type Telegram struct {
	ChatId   string `yaml:"chatid"`
	BotToken string `yaml:"bottoken"`
}

type Fcm struct {
	Serverkey string `yaml:"serverkey"`
}

type Sms struct {
	UserKey string `yaml:"userkey"`
	PassKey string `yaml:"passkey"`
}

type Ppob struct {
	Username    string `yaml:"username"`
	PpobKey     string `yaml:"ppobkey"`
	PrepaidHost string `yaml:"prepaidhost"`
	PostpaidUrl string `yaml:"postpaidurl"`
}

type Inveli struct {
	LoanProductID string `yaml:"loanproductid"`
	BupdaAccount  string `yaml:"bupdaaccount"`
	InveliAPI     string `yaml:"inveliapi"`
}

type ApplicationConfiguration struct {
	Application   Application
	Webserver     Webserver
	Database      Database
	Jwt           Jwt
	Timezone      Timezone
	Log           Log
	IpaymuPayment IpaymuPayment
	Email         Email
	Telegram      Telegram
	Fcm           Fcm
	Sms           Sms
	Ppob          Ppob
	Inveli        Inveli
}

var lock = sync.Mutex{}
var applicationConfiguration *ApplicationConfiguration

func GetConfig() *ApplicationConfiguration {
	lock.Lock()
	defer lock.Unlock()

	if applicationConfiguration == nil {
		applicationConfiguration = initConfig()
	}

	return applicationConfiguration
}

func initConfig() *ApplicationConfiguration {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var finalConfig ApplicationConfiguration
	err := viper.Unmarshal(&finalConfig)
	if err != nil {
		panic(err)
	}
	return &finalConfig
}
