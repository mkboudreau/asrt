package commands

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func OsSignalShutdown(doBeforeShutdown func(), shutdownDelay int) {
	sigc := make(chan os.Signal, 2)

	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func(s <-chan os.Signal, fn func(), delay int) {
		str := <-s
		log.Println("Received signal:", str.String())
		fn()

		log.Printf("Shutting down after %v seconds\n", delay)
		time.Sleep(time.Duration(delay) * time.Second)

		os.Exit(0)
	}(sigc, doBeforeShutdown, shutdownDelay)
}
