package initializers

import (
	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/db_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GrpcClient desc.DatabaseServiceClient

func ConnectGRPC(url string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	GrpcClient = desc.NewDatabaseServiceClient(conn)
	return conn, nil
}
