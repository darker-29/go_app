package main

// import "github.com/gin-gonic/gin"
// import "encoding/json"

import (
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"src/db_conn"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/todo", func(w http.ResponseWriter, q *http.Request) {
		m := q.Method

		if m == "GET" {
			todoList := db_conn.GetAllRecord()
			returnTemplateWithMultiStruct("./templates/index.html", todoList, w)
		}

		if m == "POST" {
			sm := q.PostFormValue("_method")

			if sm == "PUT" {
				db_conn.UpdateById(q)
			} else if sm == "DELETE" {
				db_conn.DeleteById(q)
			} else {
				db_conn.CreateTodo(q)
			}

			http.Redirect(w, q, "/todo", 301)
		}
	})

	mux.HandleFunc("/new", func(w http.ResponseWriter, q *http.Request) {
		returnTemplateWithNil("./templates/new.html", w)
	})

	mux.HandleFunc("/edit", func(w http.ResponseWriter, q *http.Request) {
		todo := db_conn.GetTitleById(q)
		returnTemplateWithMonoStruct("./templates/edit.html", todo, w)
	})

	http.ListenAndServe(":3333", mux)
}

func returnTemplateWithMonoStruct(temp string, todo db_conn.Todo, w http.ResponseWriter) {
	t, err := template.ParseFiles(temp)
	if err != nil {
		panic(err.Error())
	}
	if err := t.Execute(w, todo); err != nil {
		panic(err.Error())
	}
}

func returnTemplateWithMultiStruct(temp string, todos []db_conn.Todo, w http.ResponseWriter) {
	t, err := template.ParseFiles(temp)
	if err != nil {
		panic(err.Error())
	}
	if err := t.Execute(w, todos); err != nil {
		panic(err.Error())
	}
}

func returnTemplateWithNil(temp string, w http.ResponseWriter) {
	t, err := template.ParseFiles(temp)
	if err != nil {
		panic(err.Error())
	}
	if err := t.Execute(w, nil); err != nil {
		panic(err.Error())
	}
}
