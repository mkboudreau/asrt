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
	"github.com/mkboudreau/asrt/execution"
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

	executor := execution.NewExecutor(c.AggregateOutput, c.FailuresOnly, c.StateChangeOnly, c.ResultFormatter(), c.Writer(), c.Workers)

	asrt := NewAsrtHandler(c, executor)
	asrt.refreshServerCache()
	go asrt.loopServerCacheRefresh()

	http.Handle("/data", asrt)
	http.HandleFunc("/", serveStaticWebFiles)

	fmt.Println("Listening on port:", ctx.String("port"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", ctx.String("port")), nil))
}

type AsrtHandler struct {
	configuration *config.Configuration
	executor      *execution.Executor
	cachedContent []byte
	buffer        *bytes.Buffer
	mutex         *sync.RWMutex
}

func NewAsrtHandler(cfg *config.Configuration, executor *execution.Executor) *AsrtHandler {
	return &AsrtHandler{configuration: cfg, executor: executor, mutex: &sync.RWMutex{}}
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
	asrt.executor.Execute(asrt.configuration.Targets)
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
