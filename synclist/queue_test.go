package synclist

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewSyncList(t *testing.T) {
	sl := NewSyncList(context.Background(), &Options{
		Debug:                true,
		FlushLength:          10,
		MaxFlushIntervalTime: 60,
		MinFlushIntervalTime: 30,
		UseRandIntervalTime:  true,
		Topic:                "chat",
		StopSignalTopic:      "_stop",
		Handler: func(s *SyncList) error {
			return nil
		},
	})
	assert.NotNil(t, sl)
	assert.NotNil(t, sl.Lock)
	assert.NotNil(t, sl.Data)
	assert.LessOrEqualf(t, sl.flushIntervalTime, int64(90), "错误")
	assert.NotNil(t, sl.handler)
	assert.NotNil(t, sl.once)
}

func TestSyncList_Close(t *testing.T) {
	type dataItem struct {
		Id   int
		Name string
		Age  int
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler := func(s *SyncList) error {
		fmt.Println("长度", s.Data.Len())
		for el := s.Data.Front(); el != nil; {
			fmt.Println(el.Value)
			next := el.Next()
			s.Data.Remove(el)
			el = next
		}
		fmt.Println("长度", s.Data.Len())
		return nil
	}

	sl := NewSyncList(ctx, &Options{
		Debug:                true,
		FlushLength:          3,
		MaxFlushIntervalTime: 5,
		MinFlushIntervalTime: 2,
		UseRandIntervalTime:  true,
		Topic:                "chat",
		StopSignalTopic:      "_stop",
		Handler:              handler,
	})

	sl.StartJob()
	sl.StartJob()
	sl.StartJob()

	for i := 0; i < 30; i++ {
		sl.Lock.Lock()
		sl.Data.PushBack(dataItem{
			Id:   i,
			Name: fmt.Sprintf("name_%d", i),
			Age:  i + 10,
		})
		fmt.Println(i)
		sl.Lock.Unlock()
		if i%3 == 0 {
			time.Sleep(time.Second * 3)
		}
	}
	fmt.Println("结束。。。")
	sl.StopJob()

	fmt.Println("测试结束")
	assert.Nil(t, sl.Data, "sl.Data should be nil")
}

func TestSyncList_Execute(t *testing.T) {
	type dataItem struct {
		Id   int
		Name string
		Age  int
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	handler := func(s *SyncList) error {
		fmt.Println("长度", s.Data.Len())
		for el := s.Data.Front(); el != nil; {
			fmt.Println(el.Value)
			next := el.Next()
			s.Data.Remove(el)
			el = next
		}
		fmt.Println("长度", s.Data.Len())
		return nil
	}

	sl := NewSyncList(ctx, &Options{
		Debug:                true,
		FlushLength:          3,
		MaxFlushIntervalTime: 5,
		MinFlushIntervalTime: 2,
		UseRandIntervalTime:  true,
		Topic:                "chat",
		StopSignalTopic:      "_stop",
		Handler:              handler,
	})

	sl.StartJob()
	sl.StartJob()
	sl.StartJob()

	for i := 0; i < 30; i++ {
		sl.Lock.Lock()
		sl.Data.PushBack(dataItem{
			Id:   i,
			Name: fmt.Sprintf("name_%d", i),
			Age:  i + 10,
		})
		fmt.Println(i)
		sl.Lock.Unlock()
		if i%3 == 0 {
			time.Sleep(time.Second * 3)
		}
		if i == 11 {
			cancel()
			break
		}
	}
	fmt.Println("结束。。。")
	sl.StopJob()

	fmt.Println("测试结束")
	assert.Nil(t, sl.Data, "sl.Data should be nil")
}

func TestNewRandIntervalTime(t *testing.T) {
	var min int64 = 10
	var div int64 = 60
	a := NewRandIntervalTime(min, div)
	assert.GreaterOrEqual(t, a, min, "错误")
	assert.LessOrEqualf(t, a, min+div, "错误")

	a = NewRandIntervalTime(min, div)
	assert.GreaterOrEqual(t, a, min, "错误")
	assert.LessOrEqualf(t, a, min+div, "错误")

	a = NewRandIntervalTime(min, div)
	assert.GreaterOrEqual(t, a, min, "错误")
	assert.LessOrEqualf(t, a, min+div, "错误")
}
