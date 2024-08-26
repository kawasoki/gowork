package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Job2 struct {
	InitData   int64
	ResultData int64
}

var total int64

func countData(data []*Job2) {
	for _, datum := range data {
		atomic.AddInt64(&total, datum.InitData)
	}
}

func rWorker1(wg *sync.WaitGroup, id int, jobs <-chan []*Job2) {
	defer wg.Done()
	for job := range jobs {
		fmt.Println("worker", id, "processing job")
		countData(job)
	}
}
func DoWork() {
	jobs := make(chan []*Job2, 100)
	var wg sync.WaitGroup
	wg.Add(3)
	for w := 1; w <= 3; w++ {
		go rWorker1(&wg, w, jobs)
	}
	var j int64
	var jobS []*Job2
	for j = 1; j <= 1000000; j++ {
		jobS = append(jobS, &Job2{InitData: j})
		if len(jobS) == 100 {
			jobs <- jobS
			//jobS = jobS[:0]
			jobS = make([]*Job2, 0)
		}
	}
	fmt.Println(1111)
	close(jobs)
	wg.Wait()
	fmt.Println("结束", total)
}
