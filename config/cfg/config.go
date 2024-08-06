package cfg

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DataBase struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	}
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	}
	Server struct {
		Port        string `yaml:"port"`
		Environment string `yaml:"environment"`
	}
}

var Cfg = Config{}

func InitCfg() error {

	data, err := ioutil.ReadFile("/library/config/cfg/config.yaml")

	if err != nil {
		log.Println("Error read config: ", err)
		return err
	}

	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		log.Println("Error unmarshal config: ", err)
		return err
	}
	log.Println("Init config")
	return nil
}
