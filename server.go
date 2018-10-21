package main

import (
	"flag"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"gopkg.in/h2non/filetype.v1"
	"gopkg.in/h2non/filetype.v1/types"
	"log"
	"strings"
	"time"
)

var (
	addr     = flag.String("addr", ":5000", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")

	githubPattern = "https://raw.githubusercontent.com/%s/%s/%s/%s"
	gitlabPattern = "https://gitlab.com/%s/%s/raw/%s/%s"
)

func main() {
	flag.Parse()

	types.Add(types.NewType("css", "text/css"))
	types.Add(types.NewType("json", "application/js"))
	types.Add(types.NewType("js", "application/json"))

	router := fasthttprouter.New()
	router.GET("/gitlab/:owner/:repo/:branch/*file", serve(gitlabPattern))
	router.GET("/github/:owner/:repo/:branch/*file", serve(githubPattern))
	router.GET("/", func(ctx *fasthttp.RequestCtx) {
		ctx.SendFile("docs.html")
	})

	handler := router.Handler
	if *compress {
		handler = fasthttp.CompressHandler(handler)
	}

	fmt.Println(fmt.Sprintf("TypedGit listening on %s!", *addr))
	if err := fasthttp.ListenAndServe(*addr, handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func serve(basepath string) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		status, body, err := fasthttp.GetTimeout(nil, fmt.Sprintf(basepath, ctx.UserValue("owner"), ctx.UserValue("repo"), ctx.UserValue("branch"), ctx.UserValue("file")), 5*time.Second)
		if status != 200 || err != nil {
			ctx.Error("Unable to download the file", 500)
			return
		}

		file := strings.Split(fmt.Sprintf("%s", ctx.UserValue("file")), ".")
		mime := filetype.GetType(file[len(file)-1]).MIME.Value
		if mime == "" {
			mime = "text/plain"
		}

		ctx.SetContentType(mime)
		ctx.SetBody(body)
	}
}
