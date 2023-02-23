package main

import (
	pb "FileTrasfer/Server/proto"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

// Buffer バッファサイズ -> 1Mbyte
const Buffer = 1024 * 1024

type uploadServer struct {
	pb.UnsafeFileTransferServer
}

func (u uploadServer) UploadFile(stream pb.FileTransfer_UploadFileServer) error {

	first, _ := stream.Recv()
	wrtfile, err := os.Create(first.Filename)
	if err != nil {
		return err
	}

	defer wrtfile.Close()

	wrtfile.Write(first.Filedata[0:first.Size])
	fmt.Printf("Size:%d\n", first.Size)

	for {
		s, err := stream.Recv()
		// データ終了確認
		if err == io.EOF {
			break
		}
		// エラー発生
		if err != nil {
			return err
		}
		wrtfile.Write(s.Filedata[0:s.Size])
		fmt.Printf("Size:%d\n", s.Size)
	}

	// 終了処理
	if err := stream.SendAndClose(&pb.UploadResponse{
		Result: true,
	}); err != nil {
		return err
	} else {
		return nil
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFileTransferServer(s, &uploadServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
