package user

import (
	"database/sql"
	"time"
)

type SignUpUser struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PassWord  string    `json:"password"`
	CreatedAt time.Time `json:"createdat"`
	UpdatedAt time.Time `json:"updatedat"`
}

//Userの構造体を取得
func NewSignUpUser() *SignUpUser {
	return new(SignUpUser)
}

func (user *SignUpUser) SingUp(name string, email string, password string, db *sql.DB) error {
	nowTime := time.Now()
	//DBに送るuser情報
	user.Name = name
	user.Email = email
	user.PassWord = password
	user.CreatedAt = nowTime
	user.UpdatedAt = nowTime

	//DBにuser情報を登録
	_, err := db.Exec("INSERT INTO users (Name,Email,PassWord,CreatedAt,UpdatedAt) VALUES(?,?,?,?,?)", user.Name, user.Email, user.PassWord, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	return err

}
