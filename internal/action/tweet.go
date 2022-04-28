package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kurrik/oauth1a"
)

const (
	TweetEndpoint = "https://api.twitter.com/2/tweets"
)

type TwitterConfig struct {
	service    *oauth1a.Service
	userConfig *oauth1a.UserConfig
}

type Tweet struct {
	Text string         `json:"text"`
	cfg  *TwitterConfig `json:"-"`
}

func (cfg *TwitterConfig) NewAction(text string) Action {
	return &Tweet{
		Text: text,
		cfg:  cfg,
	}
}

func (t *Tweet) Execute() error {
	httpClient := new(http.Client)
	payload, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("error marshaling tweet payload: %s", err)
	}
	httpRequest, err := http.NewRequest("POST", TweetEndpoint, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error building tweet request: %s", err)
	}
	httpRequest.Header.Add("Content-type", "application/json")
	t.cfg.service.Sign(httpRequest, t.cfg.userConfig)
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("error sending tweet request: %s", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading tweet response body: %s", err)
		body = []byte("")
	}
	if resp.StatusCode != 201 {
		return fmt.Errorf("error creating tweet: %s: %s", resp.Status, string(body))
	}
	return nil
}

func NewTwitterConfig(consumerKey, consumerSecret, accessToken, accessTokenSecret string) *TwitterConfig {

	service := &oauth1a.Service{
		RequestURL:   "https://api.twitter.com/oauth/request_token",
		AuthorizeURL: "https://api.twitter.com/oauth/authorize",
		AccessURL:    "https://api.twitter.com/oauth/access_token",
		ClientConfig: &oauth1a.ClientConfig{
			ConsumerKey:    consumerKey,
			ConsumerSecret: consumerSecret,
			CallbackURL:    "oob",
		},
		Signer: new(oauth1a.HmacSha1Signer),
	}
	user := oauth1a.NewAuthorizedConfig(accessToken, accessTokenSecret)

	return &TwitterConfig{
		service:    service,
		userConfig: user,
	}
}
