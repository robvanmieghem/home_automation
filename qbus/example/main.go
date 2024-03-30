package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/robvanmieghem/home_automation/qbus"
)

func main() {
	url := flag.String("url", "", "http://controllerip:port")
	user := flag.String("user", "", "username")
	pass := flag.String("password", "", "password")
	flag.Parse()
	ctdClient := qbus.NewClient(*url, *user, *pass)
	err := ctdClient.Login()
	if err != nil {
		log.Fatal(err)
	}
	groups, err := ctdClient.GetGroups()
	if err != nil {
		log.Fatal(err)
	}
	jsonGroups, _ := json.Marshal(groups)
	fmt.Println(string(jsonGroups))
}
