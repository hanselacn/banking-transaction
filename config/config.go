package cfg

type Config struct {
	DB     Database
	Server Server
	Worker Worker
}

type Database struct {
	Driver   string
	Name     string
	Host     string
	User     string
	Password string
}

type Server struct {
	Address string
	Port    string
	TLS     string
}

type Worker struct {
	PayoutInterval string
	PayoutTimeUnit string
}
