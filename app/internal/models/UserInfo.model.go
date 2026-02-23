package models

import desc "github.com/XenonPPG/KRS_CONTRACTS/gen/user_v1"

type UserInfo struct {
	Id   int64         `json:"user_id"`
	Role desc.UserRole `json:"role"`
}
