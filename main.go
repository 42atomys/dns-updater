package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const WebhookURL = "https://canary.discord.com/api/webhooks/858713820200566804/NbsedN-G2yzbtM2vM9TyKXODYe4Jw0HVtC_AcZxPk9yTsqA5LhBsAxsBo23SYFJ0hKmK"

type Content struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

func main() {

	log.Println("Getting current IP")
	resp, err := http.Get("https://ifconfig.co/ip")
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	c := Content{
		Content:  fmt.Sprintf("IP: %s", string(bodyBytes)),
		Username: "DNS Updater",
	}

	var jsonData []byte
	jsonData, err = json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Post to Discord")
	_, err = http.Post(WebhookURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Done !")
}

//  https://discord.com/api/webhooks/858713820200566804/NbsedN-G2yzbtM2vM9TyKXODYe4Jw0HVtC_AcZxPk9yTsqA5LhBsAxsBo23SYFJ0hKmK
