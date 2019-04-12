package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/pem"
	"io"
	"log"
	"net"

	config "github.com/niktrix/bufferframework/server/config"

	pb "github.com/niktrix/bufferframework"

	"google.golang.org/grpc"
)

type server struct{}

func (s server) Max(srv pb.Find_MaxServer) error {

	var max int32
	ctx := srv.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		PEMBlock, _ := pem.Decode(req.Key)
		if PEMBlock == nil {
			log.Fatal("Could not parse Private Key PEM")
		}

		hashed := sha256.Sum256([]byte(string(req.Num)))

		var pk rsa.PublicKey
		asn1.Unmarshal(PEMBlock.Bytes, &pk)

		err = rsa.VerifyPKCS1v15(&pk, crypto.SHA256, hashed[:], req.SignedData)
		if err != nil {
			log.Println("data is not signed by the same suer")
		} else {
			if req.Num <= max {
				continue
			}
		}
		max = req.Num
		resp := pb.Res{Result: max}
		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
	}
}

func main() {
	// create listiner
	configuration := config.Config()
	lis, err := net.Listen("tcp", configuration.Server)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFindServer(s, server{})

	log.Println("Running server...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
