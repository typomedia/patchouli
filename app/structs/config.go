package structs

type Config struct {
	Company  string `json:"company"`
	Email    string `json:"email"`
	Interval int    `json:"interval"`
}

type Smtp struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
