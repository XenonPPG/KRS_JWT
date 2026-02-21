package models

type GetAllItemsRequest struct {
	Limit  int32 `query:"limit" json:"limit"`
	Offset int32 `query:"offset" json:"offset"`
}
