package pubsub

import (
	"context"
	"sync"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/util"
)

// In-memory
// Buffer channel as queue
// Transmission of messages with specific topic to all subscribers within a group
type localPubSub struct {
	id           string
	messageQueue chan *Message
	mapTopic     map[Topic][]chan *Message
	locker       *sync.RWMutex
	logger       appctx.Logger
}

func NewLocalPubSub(id string) *localPubSub {
	return &localPubSub{
		id:           id,
		messageQueue: make(chan *Message, 10000),
		mapTopic:     make(map[Topic][]chan *Message),
		locker:       new(sync.RWMutex),
	}
}

func (ps *localPubSub) Publish(ctx context.Context, topic Topic, msg *Message) error {
	msg.SetTopic(topic)

	go func() {
		defer util.Recovery()

		ps.messageQueue <- msg
		ps.logger.Info("New message published :", msg.String())
	}()

	return nil
}

func (ps *localPubSub) Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, unsubscribe func()) {
	c := make(chan *Message)

	ps.locker.Lock()

	val, ok := ps.mapTopic[topic]
	if ok {
		val = append(ps.mapTopic[topic], c)
		ps.mapTopic[topic] = val
	} else {
		ps.mapTopic[topic] = []chan *Message{c}
	}

	ps.locker.Unlock()

	return c, func() {
		ps.logger.Info("Unsubscribe :", topic)

		if chans, ok := ps.mapTopic[topic]; ok {
			for index := range chans {
				if chans[index] == c {
					chans = append(chans[:index], chans[index+1:]...)

					ps.locker.Lock()
					ps.mapTopic[topic] = chans
					ps.locker.Unlock()
				}
			}
		}
	}
}

func (ps *localPubSub) ID() string {
	return ps.id
}

func (*localPubSub) InitFlags() {
}

// Send message from message queue to subscribed topic channels
func (ps *localPubSub) Run(ac appctx.AppContext) error {
	ps.logger = ac.Logger(ps.id)

	go func() {
		defer util.Recovery()

		for {
			msg := <-ps.messageQueue
			ps.logger.Info("Message dequeue :", msg.String())

			ps.locker.RLock()

			if subs, ok := ps.mapTopic[msg.Topic()]; ok {
				for index := range subs {
					go func(c chan *Message) {
						defer util.Recovery()
						c <- msg
					}(subs[index])
				}
			}

			ps.locker.RUnlock()
		}
	}()

	return nil
}

func (*localPubSub) Stop() error {
	return nil
}
