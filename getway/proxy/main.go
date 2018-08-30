package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/cjysmat/grpc-reflect/getway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type proxymxu struct {
}

func (s *proxymxu) ServeHTTP(rep http.ResponseWriter, res *http.Request) {

	uri := s.parseURL(res.URL.Path)

	params := s.parseParams(res)

	fullMethod := fmt.Sprintf("/%v/%v", uri.getServiceName(), uri.getMethod())

	fmt.Println(fullMethod)

	opt := grpc.WithDefaultCallOptions(grpc.CallCustomCodec(proto.DefaultGRPCCodecs["application/json"]), grpc.FailFast(false))
	gconn, err := grpc.Dial("127.0.0.1:53000", opt, grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}

	md := metadata.Pairs("Key", "v")

	// 新建一个有 metadata 的 context
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	var returnHeader = metadata.MD{}
	var returnTrailer = metadata.MD{}
	opt1 := grpc.Header(&returnHeader)
	opt2 := grpc.Trailer(&returnTrailer)

	var out interface{}

	err = gconn.Invoke(ctx, fullMethod, params, &out, opt1, opt2) //grpc.FailFast(false)

	fmt.Println(returnTrailer)

	if err != nil {
		fmt.Println(err)
		return
	}

	b, err := json.Marshal(out)
	rep.Header().Set("Content-Type", "application/json")
	rep.Write(b)
	return
}

type URI struct {
	packageName string
	serviceName string
	version     string
	method      string
}

func (uri *URI) getServiceName() string {
	st := strings.Split(uri.serviceName, ".")
	serviceName := ""
	for _, v := range st {
		serviceName += strings.ToUpper(v[:1]) + v[1:]
	}
	return fmt.Sprintf("%v.%v", uri.packageName, serviceName)
}

func (uri *URI) getMethod() string {
	return strings.ToUpper(uri.method[:1]) + uri.method[1:]
}

// 解析http表单参数，最终转换为需要穿透的grpc参数
func (p *proxymxu) parseParams(req *http.Request) map[string]interface{} {
	req.ParseForm()
	// 处理传统意义上表单的参数，这里添加body内传输的json解析支持
	// 解析后的值默认追加到表单内部
	// 支持post、get、json
	params := make(map[string]interface{})
	var err error
	for key, v := range req.Form {
		var data map[string]interface{}
		// curl post -d '{"a":"100", "b":"100"}'这种形式过来的数据
		// 会被解析到req.Form的key当中，这时候value是空值
		err = json.Unmarshal([]byte(key), &data)
		if err == nil {
			for kk, vv := range data {
				params[kk] = vv
			}
		} else {
			//常规的表单数据
			if len(v) > 0 {
				params[key] = v[0]
			} else {
				params[key] = ""
			}
		}
	}
	// 如果body中有数据，尝试使用json解析
	if req.ContentLength <= 0 {
		return params
	}
	var data map[string]interface{}
	buf := make([]byte, req.ContentLength)
	req.Body.Read(buf)
	err = json.Unmarshal(buf, &data)
	if err != nil || data == nil {
		return params
	}
	for k, dv := range data {
		params[k] = dv
	}
	return params
}

func (p *proxymxu) parseURL(url string) *URI {
	// /proto/service.add/v1/sum
	st := strings.Split(url, "/")
	if len(st) < 5 {
		return nil
	}
	return &URI{
		packageName: st[1],
		serviceName: st[2],
		version:     st[3],
		method:      st[4],
	}
}

func main() {

	listener, err := net.Listen("tcp", ":8084")
	if err != nil {
		log.Fatal(err)
	}

	mxu := &proxymxu{}

	http.Serve(listener, mxu)
}
