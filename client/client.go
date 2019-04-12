package main

import (
	"context"
	"io"
	"log"

	pb "github.com/niktrix/bufferframework"
	"github.com/niktrix/bufferframework/client/config"
	"github.com/niktrix/bufferframework/crypt"

	"google.golang.org/grpc"
)

type GrpcStreamer struct {
	maxClient pb.Find_MaxClient
	ctx       context.Context
	done      chan bool
}

type StreamClient interface {
	sendRequests(sendInt []int32)
	recieve()
}

func main() {

	configuration := config.Init()
	dataArray := configuration.Data
	//connect to grpc server
	grpcConn, err := grpc.Dial(configuration.Server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error connecting server: %v", err)
	}

	client := pb.NewFindClient(grpcConn)
	stream, err := client.Max(context.Background())
	if err != nil {
		log.Fatalf("Err recieving Max stream %v", err)
	}
	done := make(chan bool)

	streamer := GrpcStreamer{maxClient: stream, ctx: stream.Context(), done: done}
	start(streamer, dataArray)

}

func (gs GrpcStreamer) sendRequests(sendInt []int32) {
	for _, i := range sendInt {
		req := generateRequest(i)
		if err := request(gs.maxClient, &req); err != nil {
			log.Fatalf("Error sending to server %v", err)
		}
		log.Println("sent ", req.Num)
	}
	if err := gs.maxClient.CloseSend(); err != nil {
		log.Println(err)
	}
}

func start(streamer GrpcStreamer, dataArray []int32) {
	go streamer.sendRequests(dataArray)
	go streamer.recieve()

	go func() {
		<-streamer.ctx.Done()
		if err := streamer.ctx.Err(); err != nil {
			log.Println(err)
		}
		close(streamer.done)
	}()

	<-streamer.done
}

func request(client pb.Find_MaxClient, req *pb.Req) error {
	return client.Send(req)
}

func (gs GrpcStreamer) recieve() {
	for {
		resp, err := gs.maxClient.Recv()
		if err == io.EOF {
			close(gs.done)
			return
		}
		if err != nil {
			log.Fatalln("Err while recieveData ", err)
		}
		log.Println("received", resp.Result)
	}
}

//generateRequest: Generate request with data provided, creates certs and sign data,
func generateRequest(data int32) pb.Req {

	//generate certificates
	key, err := crypt.GetCerts()
	if err != nil {
		log.Println("Error generating certs")
	}
	//sign Data
	var signedData []byte
	signedData, err = crypt.SignData(string(data), key)
	if err != nil {
		log.Println("Error signing data", err)
	}
	keyBytes, _ := crypt.MarshalPublicKey(key.PublicKey)
	return pb.Req{Num: data, Key: keyBytes, SignedData: signedData}
}
