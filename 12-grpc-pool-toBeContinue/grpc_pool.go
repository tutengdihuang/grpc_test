package main

import (
	"log"
	"time"

	"github.com/flyaways/pool"
	"google.golang.org/grpc"
)

func main() {
	options := &pool.Options{
		InitTargets:  []string{"127.0.0.1:8080"},
		InitCap:      5,
		MaxCap:       30,
		DialTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 60,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	//初始化连接池
	p, err := pool.NewGRPCPool(options, grpc.WithInsecure())

	if err != nil {
		log.Printf("%#v\n", err)
		return
	}

	if p == nil {
		log.Printf("p= %#v\n", p)
		return
	}

	defer p.Close()

	//动态更新服务地址
	//options.Input()<-&[]string{}

	conn, err := p.Get()
	if err != nil {
		log.Printf("%#v\n", err)
		return
	}

	defer p.Put(conn)

	//实现各种业务代码
	//conn.DoSomething()

	//打印空闲连接数
	log.Printf("len=%d\n", p.IdleCount())
}