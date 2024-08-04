package authentication

import (
	"database/sql"
	"encoding/json"
	"library/internal/db"
	"library/internal/jwt"
	"log"
	"net/http"
)

type DataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTResponse struct {
	Token string `json:"token"`
}

type User struct {
	User_id  int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Rule     string `json:"rule"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Неверный метод", http.StatusMethodNotAllowed)
		return
	}
	var user DataUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	if !validateDataUser(user.Login, user.Password) {
		http.Error(w, "Неправильный формат логина или пароля", http.StatusNotAcceptable)
		return
	}

	if !checkExistUser(user.Login, user.Password) {
		http.Error(w, "Users not exist", http.StatusNotFound)
		return
	}

	tokenString, err := jwt.CreateToken(user.Login)
	if err != nil {
		http.Error(w, "Ошибка создания JWT", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JWTResponse{Token: tokenString})

}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Неверный метод", http.StatusMethodNotAllowed)
		return
	}
	var user DataUser
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if !validateDataUser(user.Login, user.Password) {
		http.Error(w, "Неправильный формат логина или пароля", http.StatusNotAcceptable)
		return
	}

	if !checkLogin(user.Login) {
		http.Error(w, "Login is busy", http.StatusInternalServerError)
		return
	}

	tokenString, err := jwt.CreateToken(user.Login)

	if err != nil {
		log.Println("Error create JWT")
		http.Error(w, "Ошибка создания JWT", http.StatusInternalServerError)
		return
	}

	if err := addUser(user.Login, user.Password); err != nil {
		http.Error(w, "Error add user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(JWTResponse{Token: tokenString})

}

func validateDataUser(login string, password string) bool {
	if len(password) < 6 {
		return false
	}
	notValidSymbol := "\\'\"%^&*-+=#№/.,<>"
	for _, i := range login {
		for _, j := range notValidSymbol {
			if i == j {
				return false
			}
		}
	}
	return true
}

func checkLogin(login string) bool {
	check := `SELECT login FROM users WHERE login = $1`
	row := db.GetBD().QueryRow(check, login)
	var existingLogin string
	if err := row.Scan(&existingLogin); err == sql.ErrNoRows {
		return true
	} else if err != nil {
		log.Println("Error checking login")
		return false
	}
	return false
}

func addUser(login, password string) error {
	addUser := `INSERT INTO users (login, password, rule) VALUES ($1, $2, $3)`
	if _, err := db.GetBD().Exec(addUser, login, password, "default"); err != nil {
		log.Println("Error add user: ", err)
		return err
	}
	return nil

}

func checkExistUser(login, password string) bool {
	check := `SELECT user_id, login, password, rule FROM users WHERE login = $1 AND password = $2`
	row := db.GetBD().QueryRow(check, login, password)

	var u User
	if err := row.Scan(&u.User_id, &u.Login, &u.Password, &u.Rule); err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Println("Error scan row: ", err)
		return false
	}
	return true
}
