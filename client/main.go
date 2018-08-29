package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jhump/protoreflect/dynamic/grpcdynamic"

	"github.com/jhump/protoreflect/dynamic"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/jhump/protoreflect/desc"
)

func main() {
	fd, err := loadProtoset("./../proto/test.protoset")
	if err != nil {
		log.Fatal(err)
		return
	}
	sd := fd.FindService("proto.GrpcServer")

	cc, err := grpc.Dial(":53000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cc.Close()

	stub := grpcdynamic.NewStub(cc)

	doA(stub, sd)

	t1 := time.Now()
	for i := 0; i < 100000; i++ {
		doB(stub, sd)
	}
	fmt.Println(time.Now().Sub(t1).Seconds())
}

func doA(stub grpcdynamic.Stub, sd *desc.ServiceDescriptor) {

	md := sd.FindMethodByName("FuncA")
	in := md.GetInputType()
	msgm := dynamic.NewMessage(in)

	msg := []byte(`{"ID":1,"ID1":"far"}`)
	msgm.UnmarshalJSON(msg)

	rep, err := stub.InvokeRpc(context.Background(), md, msgm)
	if err != nil {
		log.Fatal(err)
		return
	}

	dm1 := rep.(*dynamic.Message)
	js, _ := dm1.MarshalJSON()

	fmt.Printf("%s\n", js)
}

func doB(stub grpcdynamic.Stub, sd *desc.ServiceDescriptor) {

	md := sd.FindMethodByName("FuncB")
	in := md.GetInputType()
	msgm := dynamic.NewMessage(in)

	msg := []byte(`{"arry":[1,2,3,4,5,6,7,8]`)
	msgm.UnmarshalJSON(msg)

	rep, err := stub.InvokeRpc(context.Background(), md, msgm)
	if err != nil {
		log.Fatal(err)
		return
	}

	dm1 := rep.(*dynamic.Message)
	js, _ := dm1.MarshalJSON()
	_ = js

	//fmt.Printf("%s\n", js)
}

func loadProtoset(path string) (*desc.FileDescriptor, error) {
	var fds descriptor.FileDescriptorSet
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	bb, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err = proto.Unmarshal(bb, &fds); err != nil {
		return nil, err
	}
	return desc.CreateFileDescriptorFromSet(&fds)
}
