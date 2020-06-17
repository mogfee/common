package db_mysql

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"html/template"
	"os"
	"strconv"
	"strings"
)

type tableStruct struct {
	TableName  string
	StructName string
	PrimaryKey string
	Columns    []tableStructColumn
	Imports    map[string]string
}
type tableStructColumn struct {
	COLUMN_KEY     string
	COLUMN_NAME    string
	DATA_TYPE      string
	COLUMN_TYPE    string
	COLUMN_COMMENT string
	StructName     string
	StructType     string
	StructColumn   string
}

func getDbName(db2 *gorm.DB) (string, error) {
	rows, err := db.DB().Query("select database()")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var dbName string
	rows.Next()
	err = rows.Scan(&dbName)
	return dbName, err
}
func getTables(db2 *gorm.DB) ([]string, error) {
	tables := []string{}
	rows, err := db.DB().Query("show tables")
	if err != nil {
		return tables, err
	}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return tables, err
		}
		tables = append(tables, tableName)
	}
	return tables, err
}
func execTable(db *gorm.DB, dbName, tableName, path string) error {
	query := fmt.Sprintf(`select COLUMN_KEY,COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,COLUMN_COMMENT from information_schema.columns where TABLE_SCHEMA="%s" and table_name="%s" order by ORDINAL_POSITION asc`, dbName, tableName)
	res, err := db.DB().Query(query)
	if err != nil {
		return err
	}
	defer res.Close()
	table := tableStruct{
		TableName:  tableName,
		StructName: getkey(tableName),
		PrimaryKey: "",
		Columns:    nil,
		Imports:    make(map[string]string),
	}
	columns := []tableStructColumn{}
	for res.Next() {
		column := tableStructColumn{}
		if err := res.Scan(&column.COLUMN_KEY, &column.COLUMN_NAME, &column.DATA_TYPE, &column.COLUMN_TYPE, &column.COLUMN_COMMENT); err != nil {
			return err
		}
		column.StructName = getkey(column.COLUMN_NAME)
		column.StructType = getType(column.DATA_TYPE)
		column.StructColumn = fmt.Sprintf("`gorm:\"column:%s\" json:\"%s\"`", column.COLUMN_NAME, column.COLUMN_NAME)
		if column.COLUMN_KEY == "PRI" {
			column.StructColumn = fmt.Sprintf("`gorm:\"column:%s;primary_key\" json:\"%s\"`", column.COLUMN_NAME, column.COLUMN_NAME)
			table.PrimaryKey = column.StructName
		}
		switch column.StructType {
		case "time.Time":
			table.Imports[`"time"`] = `"time"`
		}
		columns = append(columns, column)
	}
	table.Columns = columns
	return parse(table, path)
}

func parse(data tableStruct, path string) error {
	if data.PrimaryKey == "" {
		fmt.Println(data.TableName)
	}
	temp := `package model
import(
 {{range .Imports}}
{{.}}
{{end}}
)
type {{.StructName}} struct { {{range .Columns}}
	{{.StructName}} {{.StructType}} {{.StructColumn}}  // {{.COLUMN_NAME}} {{.COLUMN_TYPE}} {{.COLUMN_COMMENT}} {{.COLUMN_KEY}}{{end}}
}

func (m *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}

func (m *{{.StructName}}) GetId() int64 {
	return m.{{.PrimaryKey}}
}`

	f, err := os.Create(fmt.Sprintf("%s/%s.go", path, data.TableName))
	if err != nil {
		return err
	}
	buf := bytes.NewBufferString("")
	t := template.Must(template.New("temp").Parse(temp))
	if err := t.Execute(buf, data); err != nil {
		return err
	}
	s1 := buf.String()
	s1 = strings.Replace(s1, "&#34;", "\"", -1)
	f.Write([]byte(s1))
	return nil
}

func getkey(str string) string {
	a1 := []string{}
	_, err := strconv.ParseInt(str[0:1], 10, 64)
	if err == nil {
		str = fmt.Sprintf("m%s", str)
	}
	for _, row := range strings.Split(str, "_") {
		if len(row) > 0 {
			key := fmt.Sprintf("%s%s", strings.ToUpper(row[0:1]), row[1:])
			a1 = append(a1, key)
		}
	}
	return strings.Join(a1, "")
}
func getType(str string) string {
	switch str {
	case "int":
		return "int64"
	case "text":
		return "string"
	case "varchar":
		return "string"
	case "char":
		return "string"
	case "timestamp":
		return "time.Time"
	case "bigint":
		return "int64"
	case "datetime":
		return "time.Time"
	case "tinyint":
		return "int64"
	case "float":
		return "float64"
	case "longtext":
		return "string"
	case "decimal":
		return "float64"
	case "smallint":
		return "int64"

	}
	return str
}
func GenTables(db *gorm.DB, path string) error {
	dbName, err := getDbName(db)
	if err != nil {
		return err
	}
	tables, err := getTables(db)
	if err != nil {
		return err
	}

	for _, row := range tables {
		if err := execTable(db, dbName, row, path); err != nil {
			return err
		}
	}
	return nil
}
