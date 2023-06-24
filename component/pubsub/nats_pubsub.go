package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/nats-io/nats.go"
	"github.com/spf13/pflag"
)

type natsPubSub struct {
	id         string
	url        string
	connection *nats.Conn
	logger     appctx.Logger
}

func NewNatsPubSub(id string) *natsPubSub {
	return &natsPubSub{id: id}
}

func (ps *natsPubSub) Publish(ctx context.Context, topic Topic, msg *Message) error {
	data, err := json.Marshal(msg.Data())
	if err != nil {
		ps.logger.Error(err)
		return err
	}

	if err := ps.connection.Publish(string(topic), data); err != nil {
		ps.logger.Error(err)
		return err
	}

	return nil
}

func (ps *natsPubSub) Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, unsubscribeFunc func()) {
	msgChan := make(chan *Message)

	sub, err := ps.connection.Subscribe(string(topic), func(msg *nats.Msg) {
		data := make(map[string]interface{})

		_ = json.Unmarshal(msg.Data, &data)

		newMsg := NewMessage(data)
		newMsg.SetTopic(topic)

		newMsg.SetAckFunc(func() error {
			return msg.Ack()
		})

		msgChan <- newMsg
	})

	if err != nil {
		ps.logger.Error(err)
	}

	return msgChan, func() {
		_ = sub.Unsubscribe()
	}
}

func (ps *natsPubSub) ID() string {
	return ps.id
}

func (ps *natsPubSub) InitFlags() {
	pflag.StringVar(
		&ps.url,
		"nats-url",
		nats.DefaultURL,
		fmt.Sprintf("NATS URL - Ex: %s", nats.DefaultURL),
	)
}

func (ps *natsPubSub) setupOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectWait := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectWait))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectWait)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		ps.logger.Infof("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))

	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		ps.logger.Infof("Reconnected [%s]", nc.ConnectedUrl())
	}))

	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		ps.logger.Infof("Exiting: %v", nc.LastError())
	}))

	return opts
}

func (ps *natsPubSub) Run(ac appctx.AppContext) error {
	ps.logger = ac.Logger(ps.id)

	conn, err := nats.Connect(ps.url, ps.setupOptions([]nats.Option{})...)
	if err != nil {
		ps.logger.Fatal(err)
	}

	ps.logger.Info("Connected to NATS service.")

	ps.connection = conn
	return nil
}

func (ps *natsPubSub) Stop() error {
	return nil
}
