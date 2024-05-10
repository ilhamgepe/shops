package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Operation func(ctx context.Context) error

func GracefulShutdown(ctx context.Context, ops map[string]Operation) <-chan struct{} {
	var wg sync.WaitGroup
	wait := make(chan struct{}, 1)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	<-s
	for k, v := range ops {
		wg.Add(1)
		innerKey := k
		innerV := v
		go func() {
			defer wg.Done()
			if err := innerV(ctx); err != nil {
				log.Printf("failed to shutdown %s : %v\n", innerKey, innerV)
			}
		}()
	}

	wg.Wait()
	close(s)
	close(wait)
	return wait
}
