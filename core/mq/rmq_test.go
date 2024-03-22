package mq

import (
	"context"
	"fmt"
	"testing"

	"github.com/apache/rocketmq-clients/golang"
)

func TestSend(t *testing.T) {
	ctx := context.Background()
	InitMQ(context.Background())
	SendDefaultMessage(ctx, "1111", "11113", []byte{12})
}

func TestReceive(t *testing.T) {
	ctx := context.Background()
	InitDefaultSimpleConsumer(ctx, func(m *golang.MessageView) error {
		fmt.Println(m)
		return nil
	})
	select {}
}
