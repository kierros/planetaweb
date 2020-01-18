// forms.go
package main

import (
    "html/template"
    "log"
    "net/http"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type ContactDetails struct {
    Email string
    Name string
    Message string
}


func InitDB(filepath string) *sql.DB {
    db, err := sql.Open("sqlite3", filepath)
    if err != nil { panic(err) }
    if db == nil { panic("db nil") }
    return db
}

func CreateTable(db *sql.DB) {
	// create table if not exists
	sql_table := `
   	CREATE TABLE IF NOT EXISTS subscribers(
        	Email TEXT NOT NULL PRIMARY KEY,
        	Name TEXT,
         	Message TEXT,
         	InsertedDateTime DATETIME
	);
	`
	_, err := db.Exec(sql_table)
	if err != nil { panic(err) }
}


func StoreItem(db *sql.DB, item ContactDetails) {
	sql_additem := `
	INSERT INTO subscribers(
		Email,
		Name,
		Message,
		InsertedDatetime
	) values(?, ?, ?, CURRENT_TIMESTAMP)
	`

	stmt, err := db.Prepare(sql_additem)
	if err != nil { panic(err) }
	defer stmt.Close()

	_, err2 := stmt.Exec(item.Email, item.Name, item.Message)
	if err2 != nil { panic(err2) }
}


func main() {
    
    db := InitDB("./planeta.db")
    defer db.Close()
    CreateTable(db)
    tmpl := template.Must(template.ParseFiles("forms.html"))
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) 
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }

        details := ContactDetails{
            Email: r.FormValue("email"),
            Name: r.FormValue("name"),
            Message: r.FormValue("message"),
        }
        StoreItem(db, details) 

        tmpl.Execute(w, struct{ Success bool }{true})
    })

    log.Fatal(http.ListenAndServe("localhost:8080", nil))
}