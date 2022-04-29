package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kurrik/oauth1a"
	"github.com/leoluz/gitops-actions/internal/git"
)

const (
	TweetEndpoint = "https://api.twitter.com/2/tweets"
)

type TwitterConfig struct {
	service    *oauth1a.Service
	userConfig *oauth1a.UserConfig
}

type Tweet struct {
	file *git.File
	cfg  *TwitterConfig
}

type PostTweetPayload struct {
	Text string `json:"text"`
}

func ToTweetPayload(file *git.File) (*PostTweetPayload, error) {
	content, err := ioutil.ReadFile(file.GetFullpath())
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %s", file.GetFullpath(), err)
	}
	return &PostTweetPayload{
		Text: string(content),
	}, nil
}

func (cfg *TwitterConfig) NewAction(file *git.File) Action {
	return &Tweet{
		file: file,
		cfg:  cfg,
	}
}

func (t *Tweet) Execute() error {
	if t.file.GetStatus() != git.StatusAdded {
		log.Printf("tweets are just created for new files: skipping %s file %s", t.file.GetStatus().String(), t.file.GetName())
		return nil
	}
	httpClient := new(http.Client)
	tweet, err := ToTweetPayload(t.file)
	if err != nil {
		return fmt.Errorf("error creating tweet payload: %s", err)
	}

	jsonPayload, err := json.Marshal(tweet)
	if err != nil {
		return fmt.Errorf("error marshaling tweet payload: %s", err)
	}
	httpRequest, err := http.NewRequest("POST", TweetEndpoint, bytes.NewReader(jsonPayload))
	if err != nil {
		return fmt.Errorf("error building tweet request: %s", err)
	}
	httpRequest.Header.Add("Content-type", "application/json")
	t.cfg.service.Sign(httpRequest, t.cfg.userConfig)
	log.Println("creating tweet")
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
