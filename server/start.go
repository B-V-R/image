package server

import (
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"image/pb"
	"image/storage"
	"net"
	"os"
)

func StartServer(address, location, cert, key string) {
	logger := zerolog.New(os.Stdout).With().Str("from", "server").Logger()
	diskStorage := storage.New(location, logger)
	imageServer := NewImageServer(logger, diskStorage)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		logger.Err(err).Msg(err.Error())
	}

	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		logger.Err(err).Msg(err.Error())
	}

	gs := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterImageServiceServer(gs, imageServer)
	reflection.Register(gs)
	logger.Info().Msg("Server Started")
	logger.Fatal().Msg(gs.Serve(listener).Error())
}
