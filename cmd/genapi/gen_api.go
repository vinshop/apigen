package genapi

import (
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/vinshop/apigen/cmd/genapi/models"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	schemaTable = "information_schema"

	outputFolderModel   = "models"
	outputFolderRepo    = "repositories"
	outputFolderService = "services"
)

var GenAPI = &cobra.Command{
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
		if err := g.exec(); err != nil {
			zap.S().Fatalw("Error", "error", err)
		}
		zap.S().Info("Stop command api")
	},
}

type Generator struct {
	DBUser  string
	DBPass  string
	DBHost  string
	DBPort  int32
	DBName  string
	DBTable string

	OutputFolder string

	db *sql.DB
}

func (g *Generator) connect() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", g.DBUser, g.DBPass, g.DBHost, g.DBPort, schemaTable)
	mysqlDB, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		zap.S().Fatalw("Can not connect to mysql", "error", err)
	}

	g.db = mysqlDB
}

func (g *Generator) close() {
	if err := g.db.Close(); err != nil {
		zap.S().Errorw("Error when close db connection", "error", err)
	}
}

func (g *Generator) inspect() (*models.DbModel, error) {
	q, err := g.db.Prepare("SELECT `COLUMN_NAME`, `DATA_TYPE`, `IS_NULLABLE`, `COLUMN_KEY` FROM `COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` =?")
	if err != nil {
		return nil, err
	}

	rows, err := q.Query(g.DBName, g.DBTable)
	if err != nil {
		return nil, err
	}

	m := &models.DbModel{
		ModelName:  GetCamelCase(g.DBTable),
		TableName:  g.DBTable,
		Attributes: make([]*models.DbModelAttribute, 0),
	}

	for rows.Next() {
		var columnName, dataType, isNullable, columnKey string
		err = rows.Scan(&columnName, &dataType, &isNullable, &columnKey)
		if err != nil {
			zap.S().Errorw("Error when scan rows", "error", err)
			return nil, err
		}

		attr := &models.DbModelAttribute{
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
	return m, nil
}

func (g *Generator) exec() error {
	g.connect()
	m, err := g.inspect()
	if err != nil {
		return err
	}
	if err := g.generateModel(m); err != nil {
		return err
	}
	return nil
}

//go:embed templates/model.tmpl
var modelTemplateContent string

func (g *Generator) genFile(file string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(file), 0770); err != nil {
		return nil, err
	}
	return os.Create(file)
}

func (g *Generator) generateModel(m *models.DbModel) error {
	template, err := template.New("model").Parse(modelTemplateContent)
	if err != nil {
		return err
	}

	fo, err := g.genFile(fmt.Sprintf("%v/%v/%v.go", g.OutputFolder, outputFolderModel, g.DBTable))
	if err != nil {
		return err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			zap.S().Error("Error close file", "error", err)
		}
	}()

	err = template.Execute(fo, m)
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
		return "float64"
	case "date", "datetime":
		if isNullable == "YES" {
			return "*time.Time"
		}
		return "time.Time"
	case "json":
		return "map[string]interface{}"
	default:
		return "interface{}"
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
