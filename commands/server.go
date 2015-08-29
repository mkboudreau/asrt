package commands

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/output"
)

func cmdServer(c *cli.Context) {
	config, err := getConfiguration(c)
	if err != nil {
		cli.ShowCommandHelp(c, "server")
		fmt.Println("Could not get configuration. Reason:", err)
		log.Fatalln("Exiting....")
	}

	if c.String("port") == "" {
		cli.ShowCommandHelp(c, "server")
		fmt.Println("Missing port")
		log.Fatalln("Exiting....")
	}

	asrt := NewAsrtHandler(config)
	asrt.refreshServerCache()
	go asrt.loopServerCacheRefresh()

	http.Handle("/data", asrt)
	http.HandleFunc("/", serveStaticWebFiles)

	fmt.Println("Listening on port:", c.String("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", c.String("port")), nil))
}

type AsrtHandler struct {
	config        *configuration
	cachedContent []byte
	buffer        *bytes.Buffer
	mutex         *sync.RWMutex
}

func NewAsrtHandler(c *configuration) *AsrtHandler {
	return &AsrtHandler{config: c, mutex: &sync.RWMutex{}}
}

func serveStaticWebFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("www/dist"))
	fs.ServeHTTP(w, r)
}

func (asrt *AsrtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := "text/plain"
	if asrt.config.Output == formatJSON {
		contentType = "application/json"
	}

	asrt.mutex.RLock()
	if asrt.cachedContent == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", contentType)
		w.Write(asrt.cachedContent)
	}
	asrt.mutex.RUnlock()
}

func (asrt *AsrtHandler) loopServerCacheRefresh() {
	done := make(chan struct{})
	fn := func() {
		close(done)
	}

	osSignalShutdown(fn, 5)

	ticker := time.NewTicker(asrt.config.Rate)

	for {
		select {
		case <-ticker.C:
			asrt.refreshServerCache()
		case <-done:
			return
		}
	}
}

func (asrt *AsrtHandler) refreshServerCache() {
	targetChannel := make(chan *target, asrt.config.Workers)
	resultChannel := make(chan *output.Result)

	go processTargets(targetChannel, resultChannel)

	for _, target := range asrt.config.Targets {
		targetChannel <- target
	}
	close(targetChannel)

	formatter := asrt.config.ResultFormatter()
	writer := asrt.config.WriterWithWriters(asrt)

	if asrt.config.AggregateOutput {
		processAggregatedResult(resultChannel, formatter, writer, asrt.config.FailuresOnly)
	} else {
		processEachResult(resultChannel, formatter, writer, asrt.config.FailuresOnly)
	}
}

func (asrt *AsrtHandler) Write(p []byte) (n int, err error) {
	if asrt.buffer == nil {
		asrt.buffer = new(bytes.Buffer)
	}

	return asrt.buffer.Write(p)
}

func (asrt *AsrtHandler) Close() error {
	asrt.mutex.Lock()
	asrt.cachedContent = asrt.buffer.Bytes()
	asrt.buffer = nil
	asrt.mutex.Unlock()
	return nil
}
