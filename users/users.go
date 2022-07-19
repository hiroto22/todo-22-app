package users

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"todo-22-app/auth"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type CreateUser struct {
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

func (user *CreateUser) CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "applicaiton/json")
	w.Header().Set("Access-Control-Allow-Origin", "https://todo-22-front.herokuapp.com")
	switch r.Method {
	case "OPTIONS":
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		return
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	e := godotenv.Load()
	if e != nil {
		http.Error(w, e.Error(), 500)
	}
	dbConnectionInfo := os.Getenv("DATABASE_URL")
	db, err := sql.Open("mysql", dbConnectionInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	defer db.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	var data CreateUserBody
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), 500)
	}

	name := data.Name
	email := data.Email
	password := data.PassWord

	if name == "" {
		http.Error(w, err.Error(), 500)
	}

	if email == "" {
		http.Error(w, err.Error(), 500)
	}

	if password == "" {
		http.Error(w, err.Error(), 500)
	}

	hashPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	userData := CreateUser{name, email, string(hashPassWord), time.Now(), time.Now()}

	stmt, err := db.Prepare("INSERT INTO users (Name,Email,PassWord,CreatedAt,UpdatedAt) VALUES(?,?,?,?,?)")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	_, err = stmt.Exec(userData.Name, userData.Email, userData.PassWord, time.Now(), time.Now())
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	var userd User

	err = db.QueryRow("SELECT * FROM users WHERE Email=?", email).Scan(&userd.ID, &userd.Name, &userd.Email, &userd.PassWord, &userd.CreatedAt, &userd.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	token, err := auth.CreateToken(userd.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	tokenData := tokenRes{token}

	json.NewEncoder(w).Encode(tokenData)

}

func NewCreateUser() *CreateUser {
	return new(CreateUser)
}
