package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type Obj struct {
	A int
	B string
}

func init() {
	//自定义struct和mp必须先在gob里面注册类型
	gob.Register(Obj{})
	gob.Register(map[string]string{})
}
func main() {
	r := gin.Default()

	store, _ := redis.NewStore(10, "tcp", "192.168.2.211:6379", "tbet999", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.GET("/ping", func(c *gin.Context) {

		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		// fmt.Println(session.Get("mysession"))
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		obj := Obj{
			A: 1,
			B: "123",
		}
		mp := map[string]string{}
		mp["a"] = "aaa"
		mp["b"] = "bbb"
		session.Set("mp", mp)
		session.Set("obj", obj)
		// session.Save()

		session.Set("count", count)
		session.Save()

		fmt.Println(session.Get("mp"))
		session.Clear()
		session.Save()
		c.JSON(200, gin.H{"count": count})
		// c.JSON(200, gin.H{
		// 	"message": "pong",
		// })
	})
	// r.Run(":3000")

	srv := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
