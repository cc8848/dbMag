package main

import (


	"net/http"
	"context"
	"time"
	"log"
)

//http://www.jianshu.com/p/0dc7596ba90a
//http://hahaya.github.io/use-protobuf-in-go/
func main() {
	go http.ListenAndServe(":8989",nil)
	ctx,cancel:=context.WithCancel(context.Background())
	go func() {
		time.Sleep(3*time.Second)
		cancel()
	}()
	log.Print(A(ctx))
	select {

	}
}


func C(ctx context.Context) string {
	select {
	case <-ctx.Done():
		return "C Done"
	}
	return ""
}

func B(ctx context.Context) string {
	ctx, _ = context.WithCancel(ctx)
	go log.Println(C(ctx))
	select {
	case <-ctx.Done():
		return "B Done"
	}
	return ""
}

func A(ctx context.Context) string {
	go log.Println(B(ctx))
	select {
	case <-ctx.Done():
		return "A Done"
	}
	return ""
}

