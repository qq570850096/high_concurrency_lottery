package conf

const DriverName = "mysql"

type DbConfig struct {
	Host      string
	Port      int
	User      string
	Pwd       string
	DataBase  string
	IsRunning bool
}

var DbMasterList = []DbConfig{
	{
		Host:      "127.0.0.1",
		Port:      3306,
		User:      "root",
		Pwd:       ")Yonglun1003",
		DataBase:  "lottery",
		IsRunning: true,
	},
}

var DbMaster DbConfig = DbMasterList[0]
