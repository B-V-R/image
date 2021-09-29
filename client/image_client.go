package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"image/pb"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ImageClient struct {
	Log     zerolog.Logger
	service pb.ImageServiceClient
}

func NewImageClient(log zerolog.Logger, client pb.ImageServiceClient) ImageClient {
	return ImageClient{log, client}
}

func (client *ImageClient) Upload(imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		client.Log.Err(err).Msg("cannot open image file")
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.service.Save(ctx)
	if err != nil {
		client.Log.Err(err).Msg("cannot upload image")
		return
	}

	req := &pb.Image{
		Data: &pb.Image_Info{
			Info: &pb.ImageInfo{
				Name:   file.Name(),
				Format: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		client.Log.Err(stream.RecvMsg(nil)).Msg("cannot send image info to server")
		return
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			client.Log.Err(err).Msg("cannot read chunk to buffer")
		}

		req := &pb.Image{
			Data: &pb.Image_Chunk{
				Chunk: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			client.Log.Err(stream.RecvMsg(nil)).Msg("cannot send chunk to server")
			return
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		client.Log.Err(err).Msg("cannot receive response")
	}

	client.Log.Debug().Msg("image uploaded with file name" + res.GetFilename())
	fmt.Println(res.GetUrl())
}

func Delete(client pb.ImageServiceClient, fileName string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	info := pb.ImageInfo{
		Name: fileName,
	}
	status, err := client.Delete(ctx, &info)

	if err != nil {
		log.Println("failed to delete file")
	}

	log.Println("Deleted ", status.GetUrl())
}
