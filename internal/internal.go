package internal

import (
	"context"
	"fmt"
	"github.com/autom8ter/api/go/api"
	"github.com/autom8ter/gosub"
	"github.com/autom8ter/gosub/driver"
	"github.com/nlopes/slack"
	"net/http"
	"os"
)

var SLASH_FUNCTION_URL = "https://us-central1-autom8ter-19.cloudfunctions.net/SlashFunction"
var PROJECT_ID = os.Getenv("PROJECT_ID")
var SLACK_SIGNING_SECRET = []byte(os.Getenv("SLACK_SIGNING_SECRET"))

type SlashSub struct {
	pubsub *driver.Client
}

func New(service string, middlewares ...driver.Middleware) (*SlashSub, error) {
	provider, err := gosub.NewGoSub(PROJECT_ID)
	if err != nil {
		return nil, err
	}

	s := &SlashSub{
		pubsub: &driver.Client{
			ServiceName: service,
			Provider:    provider,
			Middleware:  middlewares,
		},
	}
	driver.SetClient(s.pubsub)
	return s, nil

}

func (s *SlashSub) Client() *driver.Client {
	return s.pubsub
}

func (s *SlashSub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler()(w, r)
}

func (s *SlashSub) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			api.Util.Entry().Errorf("[ERROR] Invalid method: %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		cmd, err := slack.SlashCommandParse(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slashcmd := &api.SlashCommand{
			Token:          cmd.Token,
			TeamId:         cmd.TeamID,
			TeamDomain:     cmd.TeamDomain,
			EnterpriseId:   cmd.EnterpriseID,
			EnterpriseName: cmd.EnterpriseName,
			ChannelId:      cmd.ChannelID,
			ChannelName:    cmd.ChannelName,
			UserId:         cmd.UserID,
			UserName:       cmd.UserName,
			Command:        cmd.Command,
			Text:           cmd.Text,
			ResponseUrl:    cmd.ResponseURL,
			TriggerId:      cmd.TriggerID,
		}

		attrs := make(map[string]string)
		for _, v := range r.Cookies() {
			attrs[v.Name] = v.Value
		}

		res := driver.Publish(context.Background(), slashcmd.Command, slashcmd)
		if res.Err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s\n", "command parsed and sent to backend for processing ✔")
	}
}
