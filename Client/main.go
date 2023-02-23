package main

import (
	pb "FileTransfer/Client/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"os"
	"time"
)

const Buffer = 1024

func main() {
	// 10 バイト分のバッファを用意
	buf := make([]byte, Buffer)

	conn, err := grpc.Dial("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewFileTransferClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	cli, err := c.UploadFile(ctx)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	filename := "test.xlsx"
	r, _ := os.Open(filename)
	defer r.Close()

	for {
		count, err := r.Read(buf)
		if err == io.EOF {
			break // Reached to EOF
		}
		if err != nil {
			log.Fatal(err)
		}
		cli.Send(&pb.UploadFileStruct{
			Filename: filename,
			Size:     int32(count),
			Filedata: buf,
		})
	}
	reply, err := cli.CloseAndRecv()
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	fmt.Println(reply.GetResult())

}
