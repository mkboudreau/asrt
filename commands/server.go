package commands

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/codegangsta/cli"
	"github.com/mkboudreau/asrt/config"
	"github.com/mkboudreau/asrt/output"
)

const (
	NoFailureTextInJson string = "{'message':'No Failures'}"
	NoFailureTextInText        = "No Failures"
)

func cmdServer(ctx *cli.Context) {
	c, err := config.GetConfiguration(ctx)
	if err != nil {
		cli.ShowCommandHelp(ctx, "server")
		fmt.Println("Could not get configuration. Reason:", err)
		log.Fatalln("Exiting....")
	}

	if ctx.String("port") == "" {
		cli.ShowCommandHelp(ctx, "server")
		fmt.Println("Missing port")
		log.Fatalln("Exiting....")
	}

	asrt := NewAsrtHandler(c)
	asrt.refreshServerCache()
	go asrt.loopServerCacheRefresh()

	http.Handle("/data", asrt)
	http.HandleFunc("/", serveStaticWebFiles)

	fmt.Println("Listening on port:", ctx.String("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", ctx.String("port")), nil))
}

type AsrtHandler struct {
	configuration *config.Configuration
	cachedContent []byte
	buffer        *bytes.Buffer
	mutex         *sync.RWMutex
}

func NewAsrtHandler(cfg *config.Configuration) *AsrtHandler {
	return &AsrtHandler{configuration: cfg, mutex: &sync.RWMutex{}}
}

func serveStaticWebFiles(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("www/dist"))
	fs.ServeHTTP(w, r)
}

func (asrt *AsrtHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contentType := "text/plain"
	if asrt.configuration.Output == config.FormatJSON {
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

	OsSignalShutdown(fn, 5)

	ticker := time.NewTicker(asrt.configuration.Rate)

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
	targetChannel := make(chan *config.Target, asrt.configuration.Workers)
	resultChannel := make(chan *output.Result)

	go processTargets(targetChannel, resultChannel)

	for _, target := range asrt.configuration.Targets {
		targetChannel <- target
	}
	close(targetChannel)

	formatter := asrt.configuration.ResultFormatter()
	writer := asrt.configuration.WriterWithWriters(asrt)

	if asrt.configuration.AggregateOutput {
		processAggregatedResult(resultChannel, formatter, writer, asrt.configuration.FailuresOnly)
	} else {
		processEachResult(resultChannel, formatter, writer, asrt.configuration.FailuresOnly)
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
	if asrt.buffer == nil || asrt.buffer.Len() == 0 {
		if asrt.configuration.Output == config.FormatJSON {
			asrt.cachedContent = []byte(NoFailureTextInJson)
		} else {
			asrt.cachedContent = []byte(NoFailureTextInText)
		}
	} else {
		asrt.cachedContent = asrt.buffer.Bytes()
	}
	asrt.buffer = nil
	asrt.mutex.Unlock()
	return nil
}
