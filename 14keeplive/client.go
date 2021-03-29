package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/keepalive"
)

var addrClient = flag.String("addrClient", "localhost:50052", "the address to connect to")

var kacpClient = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addrClient, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacpClient))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	fmt.Println("Performing unary request")
	go func() {
		ticker:=time.NewTicker(10*time.Second)
		for   {
			<-ticker.C
			fmt.Println("ticker")
			res, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: "keepalive demo"})
			if err != nil {
				log.Fatalf("unexpected error from UnaryEcho: %v", err)
			}
			fmt.Println("RPC response:", res)
		}

	}()
	select {} // Block forever; run with GODEBUG=http2debug=2 to observe ping frames and GOAWAYs due to idleness.
}