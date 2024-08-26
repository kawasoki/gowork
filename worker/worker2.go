package worker

import (
	"testing"
	"time"
)

const (
	MaxWorker = 100 //随便设置值
	MaxQueue  = 200 // 随便设置值
)

// 一个可以发送工作请求的缓冲 channel
var JobQueue chan Job

func init() {
	JobQueue = make(chan Job, MaxQueue)
}

type Payload struct {
	Id int
}

func (p *Payload) Do() {
	//	fmt.Println("完成", p.Id)
}

type Job struct {
	PayLoad Payload
}

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
	}
}

// Start 方法开启一个 worker 循环，监听退出 channel，可按需停止这个循环
func (w Worker) Start() {
	go func() {
		for {
			// 将当前的 worker 注册到 worker 队列中
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel:
				// 	真正业务的地方
				//	模拟操作耗时
				time.Sleep(20 * time.Millisecond)
				job.PayLoad.Do()
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) stop() {
	go func() {
		w.quit <- true
	}()
}

// 初始化操作

type Dispatcher struct {
	// 注册到 dispatcher 的 worker channel 池
	WorkerPool chan chan Job
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool}
}

func (d *Dispatcher) Run() {
	// 开始运行 n 个 worker
	for i := 0; i < MaxWorker; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:
			go func(job Job) {
				// 尝试获取一个可用的 worker job channel，阻塞直到有可用的 worker
				jobChannel := <-d.WorkerPool
				// 分发任务到 worker job channel 中
				jobChannel <- job
			}(job)
		}
	}
}

// 接收请求，把任务加入JobQueue。
func payloadHandler(id int) {
	work := Job{PayLoad: Payload{Id: id}}
	JobQueue <- work
}

func BenchmarkPayloadHandler(b *testing.B) {
	dispatch := NewDispatcher(MaxWorker)
	dispatch.Run()
	for i := 0; i < b.N; i++ {
		payloadHandler(i)
	}
}
