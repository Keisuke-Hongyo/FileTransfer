package main

import (
	pb "FileTrasfer/Server/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

type uploadServer struct {
	pb.UnsafeFileTransferServer
}

func (u uploadServer) UploadFile(stream pb.FileTransfer_UploadFileServer) error {

	first, _ := stream.Recv()
	meta := first.GetMeta()
	wrtfile, err := os.Create(meta.Filename)
	if err != nil {
		return err
	}

	defer wrtfile.Close()

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
		wrtfile.Write(s.GetData().Binarydata[0:s.GetData().Size])

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
