package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var url string

func init() {
	flag.StringVar(&url, "url", "http://localhost:8080", "remote url default http://localhost:8080")
}

func main() {
	flag.Parse()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for _ = range ticker.C {
			httpcall()
		}
	}()

	//合建chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	//阻塞直到有信号传入
	fmt.Println("启动")
	//阻塞直至有信号传入
	s := <-c
	fmt.Println("退出信号", s)

}

func httpcall() {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
