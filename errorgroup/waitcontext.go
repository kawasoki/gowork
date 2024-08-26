package errorgroup

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"log"
	"math/rand"
	"time"
)

type User struct {
	Name string
	Age  int
}

func main() {
	start := time.Now()
	g, ctx := errgroup.WithContext(context.Background())
	f := func() error {
		return gdo(ctx)
	}
	g.Go(f)
	g.Go(f)
	if err := g.Wait(); err != nil {
		log.Println(">>>>> (!!)Something goes wrong:", err.Error())
	} else {
		log.Println("Successfully done all jobs.")
	}
	log.Println(">>>>> END:", time.Now().Sub(start).Milliseconds(), "毫秒")
	for {
		time.Sleep(time.Second)
	}
}

func gdo(ctx context.Context) error {
	var errChan = make(chan error, 1)
	go func() {
		defer close(errChan)
		err := doWork(5)
		if err != nil {
			errChan <- err
		}
	}()
	select {
	case <-ctx.Done():
		log.Println("捕获错误", context.Cause(ctx))
		return context.Cause(ctx)
	case err := <-errChan:
		return err
	}
}

func doWork(count int) error {
	r := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	index := r.Intn(2*count) + 2
	for i := 0; i < count; i++ {
		if index == i {
			return errors.New("错误")
		}
		time.Sleep(time.Second)
		log.Println("workerId:", i+1)
	}
	return nil
}
