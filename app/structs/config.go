package structs

type Config struct {
	General General `json:"general"`
	Smtp    Smtp    `json:"smtp"`
}

type General struct {
	Company  string `json:"company"`
	Email    string `json:"email"`
	Interval int    `json:"interval"`
}

type Smtp struct {
	Sender   string `form:"sender" json:"sender"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
