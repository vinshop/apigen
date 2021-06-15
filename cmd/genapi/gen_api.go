package genapi

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/vinshop/apigen/cmd/genapi/models"
	"github.com/vinshop/apigen/pkg/util"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	schemaTable = "information_schema"

	modelDir      = "models"
	repoDir       = "repositories"
	serviceDir    = "services"
	controllerDir = "controllers"
)

//go:embed templates/model.tmpl
var modelTemplate string

//go:embed templates/repo.tmpl
var repoTemplate string

//go:embed templates/service.tmpl
var serviceTemplate string

//go:embed templates/controller.tmpl
var controllerTemplate string

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

	module string
	db     *sql.DB
}

func (g *Generator) exec() error {
	if err := g.getModulePath(); err != nil {
		return err
	}

	g.connect()
	model, err := g.inspect()
	if err != nil {
		return err
	}
	if err := g.generateTemplate(modelTemplate, modelDir, model); err != nil {
		return err
	}
	repo := model.Repo()

	if err := g.generateTemplate(repoTemplate, repoDir, repo); err != nil {
		return err
	}
	service := repo.Service()
	if err := g.generateTemplate(serviceTemplate, serviceDir, service); err != nil {
		return err
	}
	controller := service.Controller()
	if err := g.generateTemplate(controllerTemplate, controllerDir, controller); err != nil {
		return err
	}
	g.format()

	return nil
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

func (g *Generator) inspect() (*models.Model, error) {
	q, err := g.db.Prepare("SELECT `COLUMN_NAME`, `DATA_TYPE`, `IS_NULLABLE`, `COLUMN_KEY` FROM `COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` =?")
	if err != nil {
		zap.S().Errorw("Error when create prepare statement", "error", err)
		return nil, err
	}

	rows, err := q.Query(g.DBName, g.DBTable)
	if err != nil {
		zap.S().Errorw("Error when exec query", "error", err)
		return nil, err
	}

	m := &models.Model{
		Pkg:    g.module + "/models",
		Module: g.module,
		Name:   util.ToSingular(GetCamelCase(g.DBTable)),
		Table:  g.DBTable,
		Fields: make([]*models.Field, 0),
	}

	for rows.Next() {
		var columnName, dataType, isNullable, columnKey string
		err = rows.Scan(&columnName, &dataType, &isNullable, &columnKey)
		if err != nil {
			zap.S().Errorw("Error when scan rows", "error", err)
			return nil, err
		}

		attr := &models.Field{
			Name:       GetCamelCase(columnName),
			Type:       GetGoDataType(dataType, isNullable),
			ColumnName: columnName,
		}

		if columnKey == "PRI" {
			attr.IsPrimaryKey = true
		}
		if isNullable == "YES" {
			attr.IsNullable = true
		}
		if attr.Type == "time.Time" || attr.Type == "*time.Time" {
			m.NeedImport = true
			m.NeedImportTime = true
		}
		m.Fields = append(m.Fields, attr)
	}
	return m, nil
}

func (g *Generator) createFile(file string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(file), 0770); err != nil {
		zap.S().Errorw("Error when create file", "error", err, "file", file)
		return nil, err
	}
	return os.Create(file)
}

func (g *Generator) getModulePath() error {
	cmd := exec.Command("bash", "-c", "go mod edit -json > gomod.json")
	if err := cmd.Run(); err != nil {
		zap.S().Errorw("Error when exec command", "error", err)
		return err
	}
	defer exec.Command("bash", "-c", "rm -f gomod.json").Run()

	f, err := os.Open("gomod.json")
	if err != nil {
		zap.S().Errorw("Error when open gomod.json", "error", err)
		return err
	}
	defer f.Close()

	var mod struct {
		Module struct {
			Path string
		}
	}

	if err := json.NewDecoder(f).Decode(&mod); err != nil {
		zap.S().Errorw("Error when decode gomod.json", "error", err)
		return err
	}

	g.module = mod.Module.Path
	if g.OutputFolder != "." {
		g.module += "/" + g.OutputFolder
	}
	return nil
}

func (g *Generator) generateTemplate(layout, folder string, data interface{}) error {
	tmpl, err := template.New("tmpl").Parse(layout)
	if err != nil {
		zap.S().Errorw("Error when parse layout", "folder", folder, "error", err)
		return err
	}
	fo, err := g.createFile(fmt.Sprintf("%v/%v/%v.go", g.OutputFolder, folder, g.DBTable))
	if err != nil {
		zap.S().Errorw("Error when create file", "output", g.OutputFolder, "folder", folder, "table", g.DBTable, "error", err)
		return err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			zap.S().Error("Error close file", "error", err)
		}
	}()

	err = tmpl.Execute(fo, data)
	if err != nil {
		zap.S().Errorw("Error when exec template", "error", err)
		return err
	}
	return nil
}

func (g *Generator) format() {
	cmd := exec.Command("gofmt","-w", g.OutputFolder)
	zap.S().Infow("cmd", "cmd",cmd.String())
	if err := cmd.Run(); err != nil {
		zap.S().Warnw("Error when format output", "error", err)
	}
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
