package slashsub

import (
	"context"
	"github.com/autom8ter/api/go/api"
	"github.com/autom8ter/gosub"
	"github.com/autom8ter/gosub/driver"
	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"net/http"
)

type SlashSub struct {
	pubsub driver.Provider
}

func  (s *SlashSub) New(sub *gosub.GoSub) (*SlashSub) {
	return &SlashSub{
		pubsub: sub,
	}
}

func (s *SlashSub) HandlerFunc(ctx context.Context, topic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd,err := slack.SlashCommandParse(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		slashcmd := &api.SlashCommand{
			Token:                cmd.Token,
			TeamId:               cmd.TeamID,
			TeamDomain:           cmd.TeamDomain,
			EnterpriseId:         cmd.EnterpriseID,
			EnterpriseName:       cmd.EnterpriseName,
			ChannelId:            cmd.ChannelID,
			ChannelName:          cmd.ChannelName,
			UserId:               cmd.UserID,
			UserName:             cmd.UserName,
			Command:              cmd.Command,
			Text:                 cmd.Text,
			ResponseUrl:          cmd.ResponseURL,
			TriggerId:            cmd.TriggerID,
		}
		attrs := make(map[string]string)
		for _, v := range r.Cookies() {
			attrs[v.Name] = v.Value
		}

		res := driver.Publish(ctx, topic, slashcmd)
		if res.Err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}


func (s *SlashSub) Router(ctx context.Context, topics []string) *mux.Router {
	mux := mux.NewRouter()
	for _, t := range topics {
		mux.Handle("/"+t, s.HandlerFunc(ctx, t))
	}
	return mux
}

