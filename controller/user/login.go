package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"todo-22-app/auth"
	"todo-22-app/middleware"
	model "todo-22-app/model/user"
	"todo-22-app/view"
)

type LoginState struct {
	Email    string `json:"email"`
	PassWord string `json:"password"`
}

//loginの際に使うAPI
func Login(w http.ResponseWriter, r *http.Request) {
	//cors
	middleware.Cors(w, r)

	//requestされたデータの読み込み
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	var data LoginState
	json.Unmarshal(body, &data)
	email := data.Email

	//データベースから情報を取得
	user := model.NewUser()
	user.Login(email)

	//requestされたemail,passwordとDBの物が正しいか確認正しければtokenを返す
	err = auth.PasswordVerify(user.PassWord, data.PassWord)
	if err != nil {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	} else {
		token, err := auth.CreateToken(user.ID)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		view.CreateToken(w, token)
	}
}