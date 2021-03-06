package slashsub

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/autom8ter/api/go/api"
	"github.com/autom8ter/gosub"
	"github.com/autom8ter/gosub/driver"
	"github.com/golang/protobuf/proto"
	"github.com/nlopes/slack"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	if SERVICE == "" {
		SERVICE = "SlashService"
	}
	if TOPIC == "" {
		TOPIC = "Slash"
	}
}

var (
	PROJECT_ID           = os.Getenv("PROJECT_ID")
	SLACK_SIGNING_SECRET = []byte(os.Getenv("SLACK_SIGNING_SECRET"))
	TOPIC                = os.Getenv("TOPIC")
	SERVICE              = os.Getenv("SERVICE")
)

type HandlerFunc func(ctx context.Context, msg *proto.Message, published *api.Msg) error

func NewHandlerFunc(fn func(ctx context.Context, msg *proto.Message, published *api.Msg) error) HandlerFunc {
	return fn
}

func (h HandlerFunc) Chain(next ...HandlerFunc) {
	h = func(ctx context.Context, msg *proto.Message, published *api.Msg) error {
		if err := h(ctx, msg, published); err != nil {
			return err
		}
		for _, h := range next {
			if err := h(ctx, msg, published); err != nil {
				return err
			}
		}
		return nil
	}
}

type SlashSub struct {
	secret []byte
	pubsub *driver.Client
}

func New(middlewares ...driver.Middleware) (*SlashSub, error) {
	provider, err := gosub.NewGoSub(PROJECT_ID)
	if err != nil {
		return nil, err
	}

	s := &SlashSub{
		secret: SLACK_SIGNING_SECRET,
		pubsub: &driver.Client{
			ServiceName: SERVICE,
			Provider:    provider,
			Middleware:  middlewares,
		},
	}
	driver.SetClient(s.pubsub)
	return s, nil

}

func (s *SlashSub) Subscribe(jSON bool, handler HandlerFunc) {
	opts := driver.HandlerOptions{
		Topic:       TOPIC,
		Name:        SERVICE,
		ServiceName: SERVICE,
		Handler:     handler,
		Concurrency: 50,
		AutoAck:     true,
		JSON:        jSON,
	}
	s.pubsub.On(opts)
}

func (s *SlashSub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.ValidateRequest(r)
	handler()(w, r)
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

func (s *SlashSub) ValidateRequest(r *http.Request) bool {
	timestamp := r.Header["X-Slack-Request-Timestamp"][0]

	// Verify the timestamp is less than 5 minutes old, to avoid replay attacks.
	now := time.Now().Unix()
	messageTime, err := strconv.ParseInt(timestamp, 0, 64)
	if err != nil {
		api.Util.Entry().Errorln("[ERROR] Invalid timestamp:", timestamp)
		return false
	}
	if math.Abs(float64(now-messageTime)) > 5*60 {
		api.Util.Entry().Errorln("[ERROR] Timestamp is from > 5 minutes from now")
		return false
	}

	// Get the signature and signing version from the HTTP header.
	parts := strings.Split(r.Header["X-Slack-Signature"][0], "=")
	if parts[0] != "v0" {
		api.Util.Entry().Errorln("[ERROR] Unsupported signing version:", parts[0])
		return false
	}
	signature, err := hex.DecodeString(parts[1])
	if err != nil {
		api.Util.Entry().Errorln("[ERROR] Invalid message signature:", parts[1])
		return false
	}

	// Read the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		api.Util.Entry().Errorln("[ERROR] Can't read request body:", err)
		return false
	}

	// Generate the HMAC hash.
	prefix := fmt.Sprintf("v0:%v:", timestamp)

	hash := hmac.New(sha256.New, s.secret)
	hash.Write([]byte(prefix))
	hash.Write(body)

	// Reset the request body so it can be read again later.
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Verify our hash matches the signature.
	return hmac.Equal(hash.Sum(nil), []byte(signature))
}
