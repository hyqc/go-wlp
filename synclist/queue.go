package synclist

import (
	"container/list"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

//	期望使用方式：
// 	注册处理方法
//	goHandle := func(s *SyncList)error{
//		// TODO 不需要加锁处理Data数据
//		return nil
//	}
//	实例化：
//	client := NewSyncList(..., goHandle)
// 	开启任务(不管调用多少次只执行一次)
//  client.StartJob()
//	结束任务：
//	client.StopJob()
//

type SyncList struct {
	Ctx                  context.Context
	Debug                atomic.Bool // 打印日志
	Data                 *list.List
	Lock                 *sync.RWMutex
	wg                   *sync.WaitGroup
	ticker               *time.Ticker
	once                 *sync.Once
	flushSignal          chan string // 刷新信号
	flushStartTime       int64       // 刷新开始时间
	flushIntervalTime    int64       // 定时器执行间隔时间秒
	flushLength          int         // 刷新落地触发条件：链表的长度超过这个值时触发
	maxFlushIntervalTime int64       // 最大间隔时长秒
	minFlushIntervalTime int64       // 最小间隔时长秒
	useRandIntervalTime  bool        // 是否开启随机定时器间隔时间，true则生成 [MinFlushIntervalTime, MinFlushIntervalTime + MaxFlushIntervalTime)之间的值作为随机间隔时间
	topic                string      // 消费的消息队列的主题
	stopSignalTopic      string      // 自定义停止信号，示例：_topic
	handler              Handler
}

type Options struct {
	Debug                bool   // 打印日志
	FlushLength          int    // 刷新落地触发条件：链表的长度超过这个值时触发
	MaxFlushIntervalTime int64  // 最大间隔时长秒
	MinFlushIntervalTime int64  // 最小间隔时长秒
	UseRandIntervalTime  bool   // 是否开启随机定时器间隔时间，true则生成 [MinFlushIntervalTime, MinFlushIntervalTime + MaxFlushIntervalTime)之间的值作为随机间隔时间
	Topic                string // 消费的消息队列的主题
	StopSignalTopic      string // 自定义停止信号，示例：_topic
	Handler              Handler
}

type Handler func(s *SyncList) error

func NewSyncList(ctx context.Context, options *Options) *SyncList {
	if options == nil {
		panic(fmt.Sprintf("options need"))
	}
	if options.StopSignalTopic == options.Topic || options.Topic == "" {
		panic(fmt.Sprintf("topic and stop singal topic cannot the same or empty"))
	}
	syn := &SyncList{
		Ctx:            ctx,
		Data:           list.New(),
		Lock:           &sync.RWMutex{},
		wg:             &sync.WaitGroup{},
		once:           &sync.Once{},
		flushSignal:    make(chan string, 1),
		flushStartTime: 0,
	}
	if options.FlushLength <= 0 {
		options.FlushLength = 100
	}
	if options.UseRandIntervalTime {
		syn.flushIntervalTime = NewRandIntervalTime(options.MinFlushIntervalTime, options.MaxFlushIntervalTime)
	} else {
		syn.flushIntervalTime = options.MaxFlushIntervalTime
	}
	syn.Debug.Store(options.Debug)
	syn.flushLength = options.FlushLength
	syn.maxFlushIntervalTime = options.MaxFlushIntervalTime
	syn.minFlushIntervalTime = options.MinFlushIntervalTime
	syn.useRandIntervalTime = options.UseRandIntervalTime
	syn.topic = options.Topic
	syn.stopSignalTopic = options.StopSignalTopic // 自定义停止信号，示例：_topic
	syn.handler = options.Handler
	syn.ticker = time.NewTicker(time.Duration(syn.flushIntervalTime) * time.Second)
	fmt.Printf("ticker time is: %d\n", syn.flushIntervalTime)
	return syn
}

func NewDefaultList(options *Options) *SyncList {
	return NewSyncList(context.Background(), options)
}

// NewRandIntervalTime 生成 min <= max <= min + div 的随机数，防止多个副本实例的定时任务处理时间间隔相同，使用随机范围，降低同一时间处理任务的概率
func NewRandIntervalTime(min, add int64) int64 {
	rand.Seed(time.Now().Unix())
	return min + rand.Int63n(add)
}

func (s *SyncList) printfln(format string, a ...any) {
	if s.Debug.Load() {
		fmt.Printf(format+"\n", a...)
	}
}

func (s *SyncList) PushBack(data interface{}) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Data.PushBack(data)
}

// Execute 出发刷新，topic为消息队列的主题或名称
func (s *SyncList) Execute() {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	total := s.Data.Len()
	execute := total >= s.flushLength || time.Now().Unix() >= s.flushStartTime+s.flushIntervalTime
	if execute {
		if err := s.handler(s); err != nil {
			s.printfln("execute goHandle error, topic: %v, error: %v\n", s.topic, err)
		}
	}
}

// goHandle 开启携程进行任务处理和监听
func (s *SyncList) goHandle() {
	go func() {
		for {
			select {
			case _, ok := <-s.ticker.C:
				if ok {
					s.flushSignal <- s.topic
				}
				s.printfln("ticker trigger...")
			case topic, ok := <-s.flushSignal:
				// 结束之前执行刷新
				s.printfln("will execute goHandle, signal value topic: %v, signal status: %v\n", topic, ok)
				if ok {
					// 主动发送的信号
					s.printfln("executing...")
					s.Execute()
				}
				if topic == s.stopSignalTopic {
					s.printfln("exit goHandler...")
					s.wg.Done()
					return
				}
			}
		}
	}()
}

// StartJob 开启后台定时任务
func (s *SyncList) StartJob() {
	s.once.Do(func() {
		s.wg.Add(1)
		s.goHandle()
	})
}

// sendStopSignal 发送停止信号
func (s *SyncList) sendStopSignal() {
	s.flushSignal <- s.stopSignalTopic
}

// close 关闭资源
func (s *SyncList) close() {
	s.Data = nil
	close(s.flushSignal)
}

// StopJob 停止工作携程并释放资源
func (s *SyncList) StopJob() {
	s.ticker.Stop()
	s.sendStopSignal()
	s.wg.Wait()
	s.close()
}
