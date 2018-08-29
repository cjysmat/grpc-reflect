package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cjysmat/grpc-reflect/proto"

	"google.golang.org/grpc"
)

func main() {

	conn, _ := grpc.Dial("127.0.0.1:53000", grpc.WithInsecure())
	client := proto.NewGrpcServerClient(conn)

	t1 := time.Now()
	for i := 0; i < 100000; i++ {
		doB(client)
	}
	fmt.Println(time.Now().Sub(t1).Seconds())
}

func doB(c proto.GrpcServerClient) {

	rep, err := c.FuncB(context.Background(), &proto.FuncbRes{Arry: []int64{1, 2, 3, 4, 5, 6, 7, 8}})
	if err != nil {
		log.Fatal(err)
	}
	_ = rep
}