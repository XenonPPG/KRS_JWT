package initializers

import (
	desc "github.com/XenonPPG/KRS_CONTRACTS/gen/db_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GrpcClient desc.DatabaseServiceClient

// TODO: replace with Config later?
func ConnectGRPC(url string) error {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	if err != nil {
		return err
	}

	GrpcClient = desc.NewDatabaseServiceClient(conn)
	return err
}
