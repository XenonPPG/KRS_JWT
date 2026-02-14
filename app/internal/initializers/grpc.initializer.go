package initializers

import (
	noteDesc "github.com/XenonPPG/KRS_CONTRACTS/gen/note_v1"
	userDesc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GrpcUserService userDesc.UserServiceClient
var GrpcNoteService noteDesc.NoteServiceClient

func ConnectGRPC(url string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	GrpcUserService = userDesc.NewUserServiceClient(conn)
	GrpcNoteService = noteDesc.NewNoteServiceClient(conn)

	return conn, nil
}
