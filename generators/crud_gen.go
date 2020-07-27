// The following directive is necessary to make the package coherent:

// +build ignore

// This program generates contributors.go. It can be invoked by running
// go generate
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/heraju/mestri"

	_ "github.com/lib/pq"
)

type Model struct {
	Package   string
	ModelName string
}

type Attr struct {
	DataType  string
	ModelName string
}

func main() {
	db, err := sql.Open("postgres", mestri.PsqlInfo)
	die(err)
	defer db.Close()

	rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")

	entities := make([]string, 0)
	models := make([]Model, 0)

	for rows.Next() {
		var table_name string
		err := rows.Scan(&table_name)
		die(err)
		buildEntity(table_name, db)
		buildHandler(table_name)
		buildUsecase(table_name)
		entities = append(entities, table_name)
		models = append(models, Model{Package: table_name, ModelName: toCamelCase(table_name)})

	}
	buildApp(entities, models)
	fmt.Print(entities)
}

func buildHandler(table_name string) bool {
	dirPath := "app/" + table_name
	fileName := dirPath + "/handler.go"
	f, err := os.Create(fileName)
	die(err)
	defer f.Close()
	appTemplate := template.Must(template.ParseFiles("templates/handler.tmpl"))
	appTemplate.Execute(f, struct {
		Timestamp time.Time
		Entity    string
		ModelName string
	}{
		Timestamp: time.Now(),
		Entity:    table_name,
		ModelName: toCamelCase(table_name),
	})
	return true
}

func buildUsecase(table_name string) bool {
	dirPath := "app/" + table_name
	fileName := dirPath + "/usecase.go"
	f, err := os.Create(fileName)
	die(err)
	defer f.Close()
	appTemplate := template.Must(template.ParseFiles("templates/usecase.tmpl"))
	appTemplate.Execute(f, struct {
		Timestamp time.Time
		Entity    string
		ModelName string
	}{
		Timestamp: time.Now(),
		Entity:    table_name,
		ModelName: toCamelCase(table_name),
	})
	return true
}

func buildApp(entities []string, models []Model) bool {
	dirPath := "app/"
	fileName := dirPath + "/app.go"
	f, err := os.Create(fileName)
	die(err)
	defer f.Close()
	appTemplate := template.Must(template.ParseFiles("templates/app.tmpl"))
	appTemplate.Execute(f, struct {
		Timestamp time.Time
		Entities  []string
		Models    []Model
	}{
		Timestamp: time.Now(),
		Entities:  entities,
		Models:    models,
	})
	return true
}

func buildEntity(table_name string, db *sql.DB) bool {
	dirPath := "app/" + table_name
	err := os.Mkdir(dirPath, 0755)

	fmt.Println("Building CRUD For --- : ", table_name)
	fileName := dirPath + "/entity.go"
	f, err := os.Create(fileName)
	die(err)
	defer f.Close()

	packageTemplate := template.Must(template.ParseFiles("templates/entity.tmpl"))

	attr, err := db.Query("select column_name, data_type from information_schema.columns where table_name = $1", table_name)

	entity := make(map[string]Attr)

	for attr.Next() {
		var column_name string
		var data_type string
		var data_type_map string
		attr.Scan(&column_name, &data_type)
		switch data_type {
		case "uuid":
			data_type_map = "string"
		case "text":
			data_type_map = "string"
		case "integer":
			data_type_map = "int64"
		case "timestamp with time zone":
			data_type_map = "string"
		case "json":
			data_type_map = "string"
		}
		entity[column_name] = Attr{DataType: data_type_map, ModelName: toCamelCase(column_name)}
	}

	packageTemplate.Execute(f, struct {
		Timestamp time.Time
		Model     string
		Entity    map[string]Attr
	}{
		Timestamp: time.Now(),
		Model:     table_name,
		Entity:    entity,
	})
	buildPgRepo(table_name, entity)
	return true
}

func buildPgRepo(table_name string, attributes map[string]Attr) bool {
	dirPath := "app/" + table_name
	err := os.Mkdir(dirPath, 0755)

	fileName := dirPath + "/pgRepository.go"
	f, err := os.Create(fileName)
	die(err)
	defer f.Close()

	repoTemplate := template.Must(template.ParseFiles("templates/pgRepository.tmpl"))
	repoTemplate.Execute(f, struct {
		Timestamp  time.Time
		Entity     string
		ModelName  string
		Attributes map[string]Attr
	}{
		Timestamp:  time.Now(),
		Entity:     table_name,
		ModelName:  toCamelCase(table_name),
		Attributes: attributes,
	})

	return true
}

// Utility functions
func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var link = regexp.MustCompile("(^[A-Za-z])|_([A-Za-z])")

func toCamelCase(str string) string {
	return link.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(strings.Replace(s, "_", "", -1))
	})
}
