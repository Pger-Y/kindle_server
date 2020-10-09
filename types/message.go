package types

type MessageReq struct {
	Uid     string `json:"uid"`
	Message string `json:"message"`
}

type MessageResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
