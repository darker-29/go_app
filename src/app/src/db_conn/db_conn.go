package db_conn

import (
	"database/sql"
	"net/http"
)

type Todo struct {
	Id    int
	Title string
}

func GetAllRecord() []Todo {
	conn, err := sql.Open("mysql", "root:@/go_app")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT id,title FROM todos where deleted_at IS NULL")
	if err != nil {
		panic(err.Error())
	}

	// カラムの名前をリストで取得
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var todoList []Todo
	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.Id, &todo.Title)
		if err != nil {
			panic(err.Error())
		}

		todoList = append(todoList, Todo{
			Id:    todo.Id,
			Title: todo.Title,
		})
	}

	return todoList
}

func GetTitleById(q *http.Request) Todo {
	conn := getConnection()
	defer conn.Close()
	query := q.URL.Query().Get("id")

	rows, err := conn.Query("SELECT id,title FROM todos" + " where deleted_at IS NULL AND id =" + query)
	if err != nil {
		panic(err.Error())
	}

	// カラムの名前をリストで取得
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var todo Todo
	for rows.Next() {
		err := rows.Scan(&todo.Id, &todo.Title)
		if err != nil {
			panic(err.Error())
		}
	}

	return todo
}

func DeleteById(q *http.Request) {
	conn := getConnection()
	id := q.PostFormValue("id")

	_, err := conn.Query("UPDATE todos SET deleted_at = ? where id = ?", 1, id)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
}

func UpdateById(q *http.Request) {
	conn := getConnection()
	id := q.PostFormValue("id")
	title := q.PostFormValue("title")

	_, err := conn.Query("UPDATE todos SET title = ? where id = ?", title, id)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
}

func CreateTodo(q *http.Request) {
	conn := getConnection()
	input := q.PostFormValue("title")

	_, err := conn.Query("INSERT INTO todos (title) VALUES (?)", input)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
}

func getConnection() *sql.DB {
	conn, err := sql.Open("mysql", "root:@/go_app")
	if err != nil {
		panic(err.Error())
	}

	return conn
}
