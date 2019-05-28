package main

import (
	"comment_bot/configs"
	"comment_bot/handlers"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func main() {
	//rename config_template.yml to config.yml
	data, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	c := configs.BotConfig{}
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		panic(err)
	}
	handler := handlers.NewHandler(c)
	handler.Start(0)
}
