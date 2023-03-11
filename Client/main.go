package main

import (
	pb "FileTransfer/Client/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"time"
)

const Buffer = 1024 * 1024

func main() {
	// 1k バイト分のバッファを用意
	buf := make([]byte, Buffer)

	conn, err := grpc.Dial("192.168.100.12:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewFileTransferClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cli, err := c.UploadFile(ctx)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}

	filename := "test.mp4"

	r, _ := os.Open(filename)
	defer r.Close()
	cli.Send(&pb.UploadFileStruct{
		File: &pb.UploadFileStruct_Meta{
			Meta: &pb.UploadFileStruct_FileMeta{Filename: filename},
		},
	},
	)

	for {
		count, err := r.Read(buf)
		if err == io.EOF {
			break // Reached to EOF
		}
		if err != nil {
			log.Fatal(err)
		}
		cli.Send(
			&pb.UploadFileStruct{
				File: &pb.UploadFileStruct_Data{
					Data: &pb.UploadFileStruct_FileData{
						Size:       int32(count),
						Binarydata: buf,
					},
				},
			},
		)
	}

	reply, err := cli.CloseAndRecv()
	if err != nil {
		if stat, ok := status.FromError(err); ok {

			fmt.Printf("code: %s\n", stat.Code())

			fmt.Printf("message: %s\n", stat.Message())

		} else {

			fmt.Println(err)

		}
		//log.Fatalf("could not greet: %v", err)
		return
	}
	fmt.Println(reply.GetResult())

}
