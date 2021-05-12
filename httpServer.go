package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"time"
)

type ServConfig struct {
	addr string
	name string
}

func main() {
	c := make(chan os.Signal, 1)

	//初始化
	s1Conf := ServConfig{
		addr: "0.0.0.0:8080",
		name: "HttpServerDemo",
	}
	s1 := http.NewServeMux()
	s1.HandleFunc("/", Index)

	//启动服务
	gs, mctx := errgroup.WithContext(context.Background())
	gs.Go(func() error {
		<-c
		fmt.Println("--signal--")
		return errors.New("signal exit")
	})

	gs.Go(func() error {
		return Run(s1Conf, s1)
	})

	gs.Go(func() error {
		msg := <-mctx.Done()
		fmt.Println("--main ctx.Done--", msg)
		return errors.New("main ctx exit")
	})

	if err := gs.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		//提示信息
		fmt.Println("Service process exited abnormally")
	}
}

func Run(conf ServConfig, handler http.Handler) error {
	s := &http.Server{
		Addr:    conf.addr,
		Handler: handler,
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		msg := <-ctx.Done()
		fmt.Println("--Server ctx.Done--", msg)
		//return s.Shutdown(context.Background())
		return errors.New("ctx exit")
	})

	g.Go(func() error {
		fmt.Println(conf.name, " is Started")
		return s.ListenAndServe()
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func Index(w http.ResponseWriter, req *http.Request) {
	time.Sleep(5 * time.Second)
	io.WriteString(w, "Home Page!\n")
}
