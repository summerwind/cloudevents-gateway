package dockerhub

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/summerwind/cloudevents-webhook-gateway/cloudevents"
)

type Webhook struct {
	Repository WebhookRepository `json:"repository"`
}

type WebhookRepository struct {
	RepoURL string `json:"repo_url"`
}

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(req *http.Request) (*cloudevents.Event, error) {
	var w Webhook

	if req.Body == nil {
		return nil, errors.New("empty payload")
	}

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	err := decoder.Decode(&w)
	if err != nil {
		return nil, err
	}

	s, err := url.Parse(w.Repository.RepoURL)
	if err != nil {
		return nil, err
	}

	ce := &cloudevents.Event{
		Type:            "com.docker.hub.push",
		Source:          *s,
		DataContentType: "application/json",
	}

	return ce, nil
}
