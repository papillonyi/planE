package main

import (
	"github.com/spf13/viper"
	"github.com/xjh22222228/ip"
	"log"
)

type AliyunConfig struct {
	Key    string
	Secret string
}

type Config struct {
	Aliyun AliyunConfig
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	var configuration Config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	log.Printf("key is %s", configuration.Aliyun.Key)
	log.Printf("secret is %s", configuration.Aliyun.Secret)
	//outboundIP, err := GetOutboundIP()
	//log.Printf("local ip is  %s", outboundIP)
}

func GetOutboundIP() (string, error) {
	return ip.V4()
}
