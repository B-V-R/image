package client

import (
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"image/pb"
	"os"
)

func NewClient(address, cert string) ImageClient {
	var opts []grpc.DialOption

	logger := zerolog.New(os.Stdout).
		With().
		Str("from", "client").
		Logger()

	creds, err := credentials.NewClientTLSFromFile(cert, "localhost")
	if err != nil {
		logger.Error().Msg(err.Error())
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))
	conn, err := grpc.Dial(address, opts...)

	if err != nil {
		logger.Error().Msg(err.Error())
	}

	client := pb.NewImageServiceClient(conn)
	return NewImageClient(logger, client)
}
