package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/cjysmat/grpc-reflect/proto"

	"google.golang.org/grpc"
)

func main() {

	//opt := grpc.WithDefaultCallOptions(grpc.CallCustomCodec(codec.DefaultGRPCCodecs["application/json"]), grpc.FailFast(false))

	conn, _ := grpc.Dial("127.0.0.1:53000", grpc.WithInsecure())
	client := proto.NewGrpcServerClient(conn)

	t1 := time.Now()
	for i := 0; i < 1; i++ {
		doB(client)
	}
	fmt.Println(time.Now().Sub(t1).Seconds())
}

func doB(c proto.GrpcServerClient) {

	md := metadata.Pairs("Key", "v")

	// 新建一个有 metadata 的 context
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	rep, err := c.FuncB(ctx, &proto.FuncbRes{Arry: []int64{1, 2, 3, 4, 5, 6, 7, 8}})
	if err != nil {
		log.Fatal(err)
	}
	_ = rep
}
