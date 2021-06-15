package modelgen

var dbType string
var dbHost string
var dbPort int32
var dbUser string
var dbPass string
var dbName string
var dbTable string
var outputFolder string

func init() {
	GenAPI.Flags().StringVarP(&dbType, "type", "", "mysql", "db type: mysql, postgres")
	GenAPI.Flags().StringVarP(&dbHost, "host", "", "localhost", "db host")
	GenAPI.Flags().Int32VarP(&dbPort, "port", "", 3306, "db port")
	GenAPI.Flags().StringVarP(&dbUser, "user", "", "root", "db user")
	GenAPI.Flags().StringVarP(&dbPass, "pass", "", "secret", "db pass")
	GenAPI.Flags().StringVarP(&dbName, "db", "", "db_name", "db name")
	GenAPI.Flags().StringVarP(&dbTable, "table", "", "table_name", "db name")
	GenAPI.Flags().StringVarP(&outputFolder, "out", "", "output", "output folder")
}
