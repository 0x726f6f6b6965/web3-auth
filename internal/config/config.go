package config

const (
	Dev = "dev"
	Pre = "pre"
	Prd = "prd"
)

type AppConfig struct {
	HttpPort     uint64   `yaml:"http_port"`
	Env          string   `yaml:"env"`
	DynamoDB     DynamoDB `yaml:"dynamodb"`
	JwtSecretKey string   `yaml:"jwt_secret_key"`
}

type DynamoDB struct {
	Table  string `yaml:"table"`
	Region string `yaml:"region"`
	Host   string `yaml:"host"`
}

func (cfg *AppConfig) IsDevEnv() bool {
	return cfg.Env == "dev"
}
