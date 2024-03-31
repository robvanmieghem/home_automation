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

const (
	ON  = 255
	OFF = 0
)

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

func (c *Client) sendCommand(command *Command) (response *Response, err error) {

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
	response = &Response{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	return
}

func (c *Client) Login() (err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return
	}
	c.httpClient = http.Client{Jar: jar}
	loginCommand := NewLoginCommand(c.User, c.Password)
	response, err := c.sendCommand(loginCommand)
	if err != nil {
		return
	}
	sessionId, err := response.GetLoginResponse()
	if err != nil {
		return
	}
	cookieURL, err := url.Parse(c.Url)
	if err != nil {
		return
	}
	jar.SetCookies(cookieURL, []*http.Cookie{{Name: "i", Value: sessionId}})
	return
}

func (c *Client) GetGroups() (groups []Group, err error) {
	command := NewGetGroupsCommand()
	response, err := c.sendCommand(command)
	if err != nil {
		return
	}
	groups, err = response.GetGroupsResponse()
	return
}

func (c *Client) SetStatus(channelID, status int) (err error) {
	command := NewSetStatusCommand(channelID, []int{status})
	response, err := c.sendCommand(command)

	if err != nil {
		return
	}
	err = response.Error()
	return
}
