package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"telegram_comment_bot/configs"
	"telegram_comment_bot/handlers"
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
	handler.Pause()

	//handler.Stop()
	//handler.Start(0)
}
