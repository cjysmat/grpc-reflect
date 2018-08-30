package main

import (
	"context"
	"fmt"
	"log"
	"net"

	codec "github.com/cjysmat/grpc-reflect/getway/proto"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	"github.com/cjysmat/grpc-reflect/proto"
)

type gserver struct {
}

func main() {
	lis, err := net.Listen("tcp", ":"+"53000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	op := grpc.CustomCodec(codec.DefaultGRPCCodecs["application/json"])
	s := grpc.NewServer(op)
	proto.RegisterGrpcServerServer(s, gserver{})
	s.Serve(lis)
}

//FuncA for ...
func (s gserver) FuncA(c context.Context, in *proto.FuncaRes) (*proto.FuncaRep, error) {
	rep := &proto.FuncaRep{}
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, nil
	}

	fmt.Println(md.Get("key"))
	fmt.Println(*in)
	rep.ID = 1
	rep.Name = "1"

	md1 := metadata.Pairs("key1", "v1")
	metadata.NewOutgoingContext(c, md1)

	return rep, nil
}

//FuncA for ...
func (s gserver) FuncB(c context.Context, in *proto.FuncbRes) (*proto.FuncbRep, error) {

	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, nil
	}

	fmt.Println(md.Get("key"))

	rep := &proto.FuncbRep{}

	for _, item := range in.Arry {
		rep.Arry = append(rep.Arry, item)
	}

	return rep, nil
}
