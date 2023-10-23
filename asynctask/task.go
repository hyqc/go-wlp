package asynctask

import (
	"encoding/json"
	"fmt"
	"go-micro.dev/v4/broker"
	"go-micro.dev/v4/logger"
	"sync"
)

// AsyncTask 异步任务
type AsyncTask struct {
	list []broker.Subscriber
	wg   *sync.WaitGroup
}

type Options struct {
	Log  logger.Logger
	Conf broker.Options
}

func NewAsyncTask(b broker.Broker) *AsyncTask {
	task := &AsyncTask{
		wg: &sync.WaitGroup{},
	}
	broker.DefaultBroker = b
	return task
}

func (a *AsyncTask) Subscribe(topic string, group string, handler broker.Handler) (err error) {
	//订阅任务
	subscribe, err := broker.DefaultBroker.Subscribe(topic, func(event broker.Event) error {
		a.wg.Add(1)
		defer a.wg.Done()
		return handler(event)
	}, broker.Queue(group))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Subscribe Topic: %v , Group: %v\n", topic, group)
	a.list = append(a.list, subscribe)
	return
}

func (a *AsyncTask) Publish(topic string, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return broker.DefaultBroker.Publish(topic, &broker.Message{
		Body: body,
	})
}

func (a *AsyncTask) Connection() {
	err := broker.DefaultBroker.Connect()
	if err != nil {
		panic(err)
	}
}

func (a *AsyncTask) Stop() {
	//关闭所有的订阅
	for _, subscribe := range a.list {
		_ = subscribe.Unsubscribe()
	}
	//等待所有再执行的任务完成再退出
	a.wg.Wait()
}
