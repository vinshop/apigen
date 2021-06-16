package genapi

var dbType string
var dbHost string
var dbPort int32
var dbUser string
var dbPass string
var dbName string
var dbTable string
var outputFolder string
var modelFolder string
var repoFolder string
var serviceFolder string
var controllerFolder string
var mapperFolder string
var dtoFolder string

func init() {
	flag := GenAPI.Flags()
	flag.StringVarP(&dbType, "type", "", "mysql", "db type: mysql, postgres")
	flag.StringVarP(&dbHost, "host", "", "localhost", "db host")
	flag.Int32VarP(&dbPort, "port", "", 3306, "db port")
	flag.StringVarP(&dbUser, "user", "", "root", "db user")
	flag.StringVarP(&dbPass, "pass", "", "secret", "db pass")
	flag.StringVarP(&dbName, "db", "", "db_name", "db name")
	flag.StringVarP(&dbTable, "table", "", "table_name", "db name")
	flag.StringVarP(&outputFolder, "out", "", "output", "output folder")

	flag.StringVarP(&modelFolder, "model", "", "models", "models folder")
	flag.StringVarP(&dtoFolder, "dto", "", "dtos", "dtos folder")
	flag.StringVarP(&repoFolder, "repo", "", "repositories", "repositories folder")
	flag.StringVarP(&serviceFolder, "service", "", "services", "services folder")
	flag.StringVarP(&controllerFolder, "controller", "", "controllers", "controllers folder")
	flag.StringVarP(&mapperFolder, "mapper", "", "mappers", "mappers folder")
}
