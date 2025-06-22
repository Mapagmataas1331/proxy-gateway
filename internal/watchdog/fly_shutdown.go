package watchdog

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ShutdownGracefully(delay time.Duration) {
	log.Printf("Shutdown scheduled in %v...", delay)
	time.Sleep(delay)
	log.Println("Shutting down now.")

	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		log.Fatalf("Failed to find process: %v", err)
	}

	if err := p.Signal(syscall.SIGTERM); err != nil {
		log.Fatalf("Failed to send SIGTERM: %v", err)
	}
}

func ListenForShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-c
		log.Printf("Received signal %v, exiting.", sig)
		os.Exit(0)
	}()
}
