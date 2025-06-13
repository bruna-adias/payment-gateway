package infra

type Configuration struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

func NewConfiguration() *Configuration {
	return &Configuration{
		DbUser:     "myuser",
		DbPassword: "mypassword",
		DbHost:     "localhost",
		DbPort:     "3307",
		DbName:     "payment_gateway",
	}
}
