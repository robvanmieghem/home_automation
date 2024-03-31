package qbus

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func NewClient(url, user, password string) *Client {
	return &Client{Url: url, User: user, Password: password}
}

type Client struct {
	Url        string
	User       string
	Password   string
	httpClient http.Client
}

type Channel struct {
	ID    uint   `json:"Chnl"`
	Name  string `json:"Nme"`
	Icon  uint   `json:"Ico"`
	Value []int  `json:"Val"`
}
type Group struct {
	Name  string    `json:"Nme"`
	Items []Channel `json:"Itms"`
}

func (c *Client) Login() (err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	c.httpClient = http.Client{Jar: jar}
	loginCommand := NewLoginCommand(c.User, c.Password)
	json_str, err := json.Marshal(loginCommand)
	if err != nil {
		return
	}
	resp, err := c.httpClient.PostForm(c.Url+"/default.aspx+r="+fmt.Sprintf("%.16f", rand.Float64()), url.Values{"strJSON": {string(json_str)}})
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Invalid response: %d", resp.StatusCode))
		return
	}

	defer resp.Body.Close()
	result := Response{}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	sessionId, err := result.GetLoginResponse()
	if err != nil {
		return
	}
	jar.SetCookies(resp.Request.URL, []*http.Cookie{{Name: "i", Value: sessionId}})
	return
}

func (c *Client) GetGroups() (groups []Group, err error) {
	command := NewGetGroupsCommand()
	json_str, err := json.Marshal(command)
	resp, err := c.httpClient.PostForm(c.Url+"/default.aspx+r="+fmt.Sprintf("%.16f", rand.Float64()), url.Values{"strJSON": {string(json_str)}})
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("Invalid response: %d", resp.StatusCode))
		return
	}

	defer resp.Body.Close()
	result := Response{}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	groups, err = result.GetGroupsResponse()
	return
}
