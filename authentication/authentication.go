package authentication

import (
	"encoding/json"
	"library/utils/jwt"
	"net/http"
)

var secretKey = []byte("DurkaVerder0304")

type DataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JWTResponse struct {
	Token string `json:"token"`
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

	tokenString, err := jwt.CreateToken(user.Login, secretKey)
	if err != nil {
		http.Error(w, "Ошибка создания JWT", http.StatusInternalServerError)
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
