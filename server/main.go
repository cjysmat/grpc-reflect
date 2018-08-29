package main

import (
	"context"
	"fmt"
	"log"
	"net"

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

	s := grpc.NewServer()
	proto.RegisterGrpcServerServer(s, gserver{})
	s.Serve(lis)
}

//FuncA for ...
func (s gserver) FuncA(c context.Context, in *proto.FuncaRes) (*proto.FuncaRep, error) {
	rep := &proto.FuncaRep{}
	fmt.Println(*in)
	rep.ID = 1
	rep.Name = "1"
	return rep, nil
}

//FuncA for ...
func (s gserver) FuncB(c context.Context, in *proto.FuncbRes) (*proto.FuncbRep, error) {
	rep := &proto.FuncbRep{}

	for _, item := range in.Arry {
		rep.Arry = append(rep.Arry, item)
	}

	return rep, nil
}
