package cmd

import (
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	schemaTable = "information_schema"

	outputFolderModel   = "models"
	outputFolderRepo    = "repositories"
	outputFolderService = "services"
)

var genAPI = &cobra.Command{
	Use:   "api",
	Short: "api",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		zap.S().Info("Start command api")

		g := &Generator{
			DBUser:  dbUser,
			DBPass:  dbPass,
			DBHost:  dbHost,
			DBPort:  dbPort,
			DBName:  dbName,
			DBTable: dbTable,

			OutputFolder: outputFolder,
		}
		g.processGenerate(args)
		zap.S().Info("Stop command api")
	},
}

type dbModel struct {
	ModelName  string
	TableName  string
	Attributes []*dbModelAttribute

	NeedImport     bool
	NeedImportTime bool
	NeedImportGorm bool
}

type dbModelAttribute struct {
	FieldName    string
	FieldType    string
	ColumnName   string
	IsPrimaryKey bool
	IsNullable   bool
}

type Generator struct {
	DBUser  string
	DBPass  string
	DBHost  string
	DBPort  int32
	DBName  string
	DBTable string

	OutputFolder string
}

func (g *Generator) processGenerate(args []string) {
	// example: user:pass@tcp(host:port)/database?param=value
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", g.DBUser, g.DBPass, g.DBHost, g.DBPort, schemaTable)
	mysqlDB, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Can not connect to mysql, detail: ", err)
	}
	defer func() {
		err = mysqlDB.Close()
		if err != nil {
			zap.S().Error(err)
		}
	}()

	stmt, err := mysqlDB.Prepare("SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_KEY FROM COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME =?")
	if err != nil {
		log.Println("Error when prepare query, detail: ", err)
		return
	}

	rows, err := stmt.Query(g.DBName, g.DBTable)
	if err != nil {
		log.Println("Error when exec query, detail: ", err)
		return
	}
	m := &dbModel{
		ModelName:  GetCamelCase(g.DBTable),
		TableName:  g.DBTable,
		Attributes: make([]*dbModelAttribute, 0),
	}

	for rows.Next() {
		var columnName, dataType, isNullable, columnKey string
		err = rows.Scan(&columnName, &dataType, &isNullable, &columnKey)
		if err != nil {
			log.Println("Error when scan rows, detail: ", err)
			return
		}
		attr := &dbModelAttribute{
			FieldName:  GetCamelCase(columnName),
			FieldType:  GetGoDataType(dataType, isNullable),
			ColumnName: columnName,
		}
		if columnKey == "PRI" {
			attr.IsPrimaryKey = true
		}
		if isNullable == "YES" {
			attr.IsNullable = true
		}
		if attr.FieldType == "time.Time" || attr.FieldType == "*time.Time" {
			m.NeedImport = true
			m.NeedImportTime = true
		}
		m.Attributes = append(m.Attributes, attr)
	}
	err = g.generateModel(m)
	if err != nil {
		zap.S().Error("Error when generateModel, detail: ", err)
	}
}

//go:embed templates/model.tmpl
var modelTemplateContent string

func (g *Generator) generateModel(m *dbModel) error {
	tmpl, err := template.New("test_model").Parse(modelTemplateContent)
	if err != nil {
		return err
	}

	// open output file
	fo, err := os.Create(fmt.Sprintf("./%v/%v/%v.go", g.OutputFolder, outputFolderModel, g.DBTable))
	if err != nil {
		return err
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			zap.S().Error("Error when exec query, detail: ", err)
			return
		}
	}()

	err = tmpl.Execute(fo, m)
	if err != nil {
		return err
	}

	return nil
}

func GetGoDataType(mysqlType, isNullable string) string {
	switch mysqlType {
	case "varchar", "longtext", "text":
		return "string"
	case "smallint", "int", "bigint", "timestamp":
		return "int64"
	case "tinyint":
		return "bool"
	case "decimal":
		return "double"
	case "date", "datetime":
		if isNullable == "YES" {
			return "*time.Time"
		}
		return "time.Time"
	default:
		return ""
	}
}

func GetCamelCase(input string) string {
	if input == "id" {
		return "ID"
	}
	output := ""
	capNext := true
	for _, v := range input {
		if v >= 'A' && v <= 'Z' {
			output += string(v)
		}
		if v >= '0' && v <= '9' {
			output += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				output += strings.ToUpper(string(v))
			} else {
				output += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}

	return output
}
