package cfg

type Config struct {
	DB DataBase
}

type DataBase struct {
	Name     string
	Address  string
	Port     string
	UserName string
	Password string
}
