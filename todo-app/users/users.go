package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PassWord  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

type CreateUserBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	PassWord string `json:"password"`
}

func (user *User) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data CreateUserBody
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatal(err)
	}

	name := data.Name
	email := data.Email
	password := data.PassWord

	hashPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	userData := User{1, name, email, string(hashPassWord), time.Now(), time.Now()}

	stmt, err := db.Prepare("INSERT INTO users (Name,Email,PassWord,CreatedAt,UpdatedAt) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	_, err2 := stmt.Exec(userData.Name, userData.Email, userData.PassWord, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err2)
	}

	w.Header().Set("Content-Type", "applicaiton/json")

	json.NewEncoder(w).Encode(userData)

}

func NewCreateUser() *User {
	return new(User)
}
