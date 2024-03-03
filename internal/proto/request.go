package proto

type GetNonceRequest struct {
	PublicAddress string `dynamodbav:"public_address" json:"public_address"`
}

type VerifySignatureRequest struct {
	PublicAddress string `dynamodbav:"public_address" json:"public_address"`
	Signature     string `json:"signature"`
}

type VerifySignatureResponse struct {
	PublicAddress string `dynamodbav:"public_address" json:"public_address"`
	Token         string `json:"token"`
}

type UpdateUserRequest struct {
	PublicAddress string   `dynamodbav:"public_address" json:"public_address"`
	User          User     `json:"user"`
	UpdateMask    []string `json:"update_mask"`
}
type UpdateUserResponse struct {
	PublicAddress string `dynamodbav:"public_address" json:"public_address"`
	User          User   `json:"user"`
	Token         string `json:"token"`
}
