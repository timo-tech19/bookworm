package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var users []User

func Signup(w http.ResponseWriter, r *http.Request) {
	var userInfo UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// create new user in database
	newUser := User{Id: len(users) + 1, Username: userInfo.Username, Email: userInfo.Email}
	users = append(users, newUser)

	token, err := createToken(newUser)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setTokenCookie(w, token)

	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var userInfo UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Invalid credentials. Please make sure all fields are valid", http.StatusBadRequest)
		return
	}

	// get database user by email
	var user User
	for _, u := range users {
		if u.Email == userInfo.Email {
			user = u
			break
		}
	}
	if user.Email == "" {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	token, err := createToken(user)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setTokenCookie(w, token)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

// Prevent users from accessing routes when unauthenicated.
func Protect(original func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err == http.ErrNoCookie {
			http.Error(w, "User is not logged in", http.StatusUnauthorized)
			return
		}

		if validateToken(token.Value) {
			original(w, r)
		} else {
			http.Error(w, "User is not logged in", http.StatusUnauthorized)
			return
		}

	}
}

// creates a jwt token for user
func createToken(user User) (string, error) {
	key := []byte("mysecretjwtauthkey")
	// creates a token struct with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
	})

	tokenString, err := token.SignedString(key)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Verifies if jwt token is valid
func validateToken(tokenString string) bool {
	key := []byte("mysecretjwtauthkey")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Could not validate auth token")
		}

		return key, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

// configures and set token cookies to be sent to client
func setTokenCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(90 * 24 * time.Hour), // 90days
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
}
