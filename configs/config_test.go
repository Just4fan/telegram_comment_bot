package configs

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"testing"
)

func TestReadConfig(t *testing.T) {
	data, err := ioutil.ReadFile("../config.yml")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	c := BotConfig{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
	log.Println(c)
}
