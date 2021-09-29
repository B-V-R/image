package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"image/pb"
	"image/storage"
	"io"
	"strings"
)

type ImageServer struct {
	Log zerolog.Logger
	pb.UnimplementedImageServiceServer
	Store storage.Storage
}

func NewImageServer(logger zerolog.Logger, store storage.Storage) pb.ImageServiceServer {
	return ImageServer{
		Store: store,
		Log:   logger,
	}
}

func (is ImageServer) Save(stream pb.ImageService_SaveServer) error {
	recv, err := stream.Recv()
	if err != nil {
		fmt.Println(err)
		return err
	}
	imageInfo := recv.GetInfo()
	imageBuf := bytes.Buffer{}
	switch strings.ToUpper(imageInfo.GetFormat()) {
	case ".PNG", ".JPG", ".GIF":
		imageSize := 0
		for {
			rec, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				fmt.Println("EOF")
				break
			}

			if err != nil {
				is.Log.Err(err).Msg("Can not receive chunk data")
				return status.Error(codes.Unknown, "Can not receive chunk data")
			}

			data := rec.GetChunk()
			imageSize += len(data)

			_, err = imageBuf.Write(data)

			if err != nil {
				is.Log.Err(err).Msg("Can not write chunk data")
				return status.Error(codes.Internal, "Can not write chunk data")
			}
		}
	default:
		return status.Error(codes.Unknown, "Un supported file format")
	}

	file, err := is.Store.Save(context.Background(), imageBuf.Bytes(), imageInfo.Name)
	if err != nil {
		is.Log.Err(err).Msg("Can not save file")
		return status.Error(codes.Internal, "Can not save file")
	}
	res := pb.Status{
		Filename: file.Name,
		Url:      fmt.Sprintf("http://localhost:8080/%s", file.URL),
	}

	err = stream.SendAndClose(&res)

	if err != nil {
		is.Log.Err(err).Msg("Can not send response")
		return status.Error(codes.Internal, "Can not send response")
	}

	is.Log.Info().Msg("File saved")

	return nil
}
func (is ImageServer) Delete(ctx context.Context, info *pb.ImageInfo) (*pb.Status, error) {
	name := info.Name

	err := is.Store.Delete(ctx, name)
	if err != nil {
		is.Log.Err(err).Msg("failed to delete file")
		return nil, status.Error(codes.Internal, "failed to delete file")
	}
	res := pb.Status{
		Filename: name,
	}
	return &res, err
}
