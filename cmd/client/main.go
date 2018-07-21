package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/johanbrandhorst/grpc-json-example/codec"
	"github.com/johanbrandhorst/grpc-json-example/insecure"
	pbExample "github.com/johanbrandhorst/grpc-json-example/proto"
)

var addr = flag.String("addr", "localhost", "The address of the server to connect to")
var port = flag.String("port", "10000", "The port to connect to")

var log grpclog.LoggerV2

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, net.JoinHostPort(*addr, *port),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(insecure.CertPool, "")),
		grpc.WithDefaultCallOptions(grpc.CallContentSubtype(codec.JSON{}.Name())),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	defer conn.Close()
	c := pbExample.NewUserServiceClient(conn)

	user := pbExample.User{Id: 1, Role: pbExample.Role_ADMIN}
	_, err = c.AddUser(ctx, &user)
	if err != nil {
		log.Fatalln("Failed to add user:", err)
	}

	srv, err := c.ListUsers(ctx, new(empty.Empty))
	if err != nil {
		log.Fatalln("Failed to list users:", err)
	}
	for {
		rcv, err := srv.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln("Failed to receive:", err)
		}
		log.Infoln("Read user:", rcv)
	}

	log.Infoln("Success!")
}
