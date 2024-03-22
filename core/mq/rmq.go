package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-clients/golang"
	"github.com/apache/rocketmq-clients/golang/credentials"
	"github.com/gogf/gf/v2/frame/g"
)

// Package main implements a simple producer to send message.
var mqProducer golang.Producer
var endpoint string
var defaultTopic string
var defaultGroupName string
var config golang.Config
var awaitDuration = time.Second * 5
var maxMessageNum int32 = 200

// invisibleDuration should > 20s
var invisibleDuration = time.Second * 20

func InitMQ(ctx context.Context) {
	_mqProducer, err := GetMqProducer(ctx)
	if err != nil {
		panic(err)
	}
	mqProducer = _mqProducer
	if err = mqProducer.Start(); err != nil {
		panic(err)
	}
}

func InitMQConfig(ctx context.Context) {
	if endpointVal, err := g.Cfg().Get(ctx, "mq.endpoint"); err != nil {
		panic("not found mq nameserver")
	} else {
		endpoint = endpointVal.String()
	}
	if defaultTopicVal, err := g.Cfg().Get(ctx, "mq.defaultTopic"); err != nil {
		panic("not found mq defaultTopic")
	} else {
		defaultTopic = defaultTopicVal.String()
	}
	if defaultGroupNameVal, err := g.Cfg().Get(ctx, "mq.defaultGroupName"); err != nil {
		panic("not found mq defaultTopic")
	} else {
		defaultGroupName = defaultGroupNameVal.String()
	}
	config = golang.Config{
		Endpoint:      endpoint,
		ConsumerGroup: defaultGroupName,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    "",
			AccessSecret: "",
		},
	}
}

func InitDefaultSimpleConsumer(ctx context.Context, fn func(m *golang.MessageView) error) {
	_mqConsumer, err := GetMqSimpleConsumer(ctx)
	if err != nil {
		panic(err)
	}
	if err := _mqConsumer.Start(); err != nil {
		panic(err)
	}
	go func() {
		for {
			mvs, err := _mqConsumer.Receive(ctx, maxMessageNum, invisibleDuration)
			if err != nil {
				fmt.Println(err)
			}
			// ack message
			for _, mv := range mvs {
				if err := fn(mv); err != nil {
					g.Log().Infof(ctx, "Receive Handle error: %v", err)
					continue
				} else {
					_mqConsumer.Ack(ctx, mv)
				}
			}
			time.Sleep(time.Second * 3)
		}
	}()
}

func GetMqSimpleConsumer(ctx context.Context) (simpleConsumer golang.SimpleConsumer, err error) {
	InitMQConfig(ctx)
	return golang.NewSimpleConsumer(&config,
		golang.WithAwaitDuration(awaitDuration),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{
			defaultTopic: golang.SUB_ALL,
		}))
}

func GetMqPushConsumer(ctx context.Context) (simpleConsumer golang.SimpleConsumer, err error) {
	InitMQConfig(ctx)
	return golang.NewSimpleConsumer(&config,
		golang.WithAwaitDuration(awaitDuration),
		golang.WithSubscriptionExpressions(map[string]*golang.FilterExpression{
			defaultTopic: golang.SUB_ALL,
		}))
}

func GetMqProducer(ctx context.Context) (mqProducer golang.Producer, err error) {
	InitMQConfig(ctx)
	return golang.NewProducer(
		&config,
		golang.WithTopics(defaultTopic),
	)
}

func SendDefaultMessage(ctx context.Context, tag string, key string, data []byte) (bool, error) {
	msg := &golang.Message{
		Topic: defaultTopic,
		Body:  data,
	}
	// set keys and tag
	msg.SetKeys(key)
	msg.SetTag(tag)
	// send message in sync
	resp, err := mqProducer.Send(ctx, msg)
	if err != nil {
		return false, err
	} else {
		g.Log().Infof(ctx, "resp: %v", resp)
	}
	return true, nil
}
