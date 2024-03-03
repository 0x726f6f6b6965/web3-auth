package proto

type User struct {
	PublicAddress string         `dynamodbav:"public_address" json:"public_address"`
	FullName      string         `dynamodbav:"full_name" json:"full_name"`
	BirthDay      string         `dynamodbav:"birth_day" json:"birth_day"`
	Nonce         string         `dynamodbav:"nonce" json:"nonce"`
	CarryId       string         `dynamodbav:"carry_id,omitempty" json:"carry_id,omitempty"`
	Addresses     []*UserAddress `dynamodbav:"addresses" json:"addresses"`

	CreatedAt int64 `dynamodbav:"created_at" json:"created_at"`
	UpdatedAt int64 `dynamodbav:"updated_at" json:"updated_at"`
}

type UserAddress struct {
	Title         string `dynamodbav:"title" json:"title"`
	StreetAddress string `dynamodbav:"street_address" json:"street_address"`
	PostalCode    string `dynamodbav:"postal_code" json:"postal_code"`
	CountryCode   int    `dynamodbav:"country_code" json:"country_code"`

	CreatedAt int64 `dynamodbav:"created_at" json:"created_at"`
	UpdatedAt int64 `dynamodbav:"updated_at" json:"updated_at"`
}

type UserToken struct {
	PublicAddress string `dynamodbav:"public_address" json:"public_address"`
	Nonce         string `dynamodbav:"nonce" json:"nonce"`
	ExpireAt      int64  `json:"expire_at"`
	CreatedAt     int64  `json:"created_at"`
}
