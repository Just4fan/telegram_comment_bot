package utils

import (
	"log"
	"testing"
)

func TestBase64Decode(t *testing.T) {
	param, err := Base64Decode("eyJtZXNzYWdlX2lkIjo4NiwiY2hhdF9pZCI6LTEwMDE0Njc4OTI4NjksImFyZWFfaWQiOjg3fQ")
	if err != nil {
		panic(err)
	}
	log.Println(string(param))
}

func TestBase64Encode(t *testing.T) {
	data := `{"message_id":86,"chat_id":-1001467892869,"area_id":87}`
	param, err := Base64Encode([]byte(data))
	if err != nil {
		panic(param)
	}
	log.Println(param)
}
