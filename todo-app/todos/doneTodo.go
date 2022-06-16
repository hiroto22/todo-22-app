package todos

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func DoneTodo(w http.ResponseWriter, r *http.Request) {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	dbConnectionInfo := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/go_todo", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"))
	db, err := sql.Open("mysql", dbConnectionInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	id := r.URL.Query().Get("id")
	isDone := r.URL.Query().Get("isDone")

	stmt, err := db.Prepare("UPDATE todos set IsDone=? WHERE ID=?")
	if err != nil {
		log.Fatal(err)
	}

	if isDone == "false" {
		_, err = stmt.Exec(true, id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err = stmt.Exec(false, id)
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Header().Set("Content-Type", "applicaiton/json")

	json.NewEncoder(w).Encode(id)

}