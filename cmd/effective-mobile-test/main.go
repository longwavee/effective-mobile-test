package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/longwavee/effective-mobile-test/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
	)
	defer stop()

	a, err := app.New()
	if err != nil {
		log.Println(fmt.Errorf("failed to init application: %w", err))
		os.Exit(1)
	}

	err = a.Start()
	if err != nil {
		log.Println(fmt.Errorf("failed to run application: %w", err))
		os.Exit(1)
	}

	<-ctx.Done()

	err = a.GracefullStop()
	if err != nil {
		log.Println(fmt.Errorf("failed to gracefull stop application: %w", err))
	}
}
