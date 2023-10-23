package asynctask

import (
	"errors"
	"fmt"
	"github.com/go-micro/plugins/v4/broker/kafka"
	"github.com/stretchr/testify/assert"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"testing"
)

func TestNewAsyncTask(t *testing.T) {
	b := kafka.NewBroker(broker.Addrs("192.168.9.100:9092"), broker.Logger(logger.NewLogger()))
	at := NewAsyncTask(b)
	assert.NotNil(t, at, "初始化错误")
}

func TestAsyncTask_Subscribe(t *testing.T) {
	b := kafka.NewBroker(broker.Addrs("192.168.9.100:9092"), broker.Logger(logger.NewLogger()))
	at := NewAsyncTask(b)
	at.Connection()
	topic := "async_task_test"
	type msg struct {
		name string
		age  int
	}
	data := "1234567890"
	fmt.Println("push data", data)
	err := at.Publish(topic, data)
	assert.Nil(t, err, err)
	err = at.Subscribe(topic, topic, func(event broker.Event) error {
		body := string(event.Message().Body)
		fmt.Println("=====", body)
		if body == data {
			return nil
		}
		return errors.New("the data cannot equal")
	})
	assert.Nil(t, err, "asynctask subscribe error", err)
	at.Stop()
	fmt.Println("pass")
}
