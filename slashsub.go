package slashsub

import (
	"context"
	"github.com/autom8ter/api/go/api"
	"github.com/autom8ter/gosub"
	"github.com/autom8ter/gosub/driver"
	"github.com/nlopes/slack"
	"net/http"
	"os"
)

type SlashSub struct {
	Project string
	pubsub  *driver.Client
}

func New(service string, middlewares ...driver.Middleware) (*SlashSub, error) {
	provider, err := gosub.NewGoSub(os.Getenv("GCP_PROJECT"))
	if err != nil {
		return nil, err
	}

	s := &SlashSub{
		Project: "",
		pubsub: &driver.Client{
			ServiceName: service,
			Provider:    provider,
			Middleware:  middlewares,
		},
	}
	driver.SetClient(s.pubsub)
	return s, nil

}

func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func SlashFunction(w http.ResponseWriter, r *http.Request) {
	s, err := New("slash", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.ServeHTTP(w, r)
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
