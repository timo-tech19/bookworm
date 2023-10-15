package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/timotech-19/bookworm/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JWTClaims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

func Signup(c *gin.Context) {
	// get user data
	var body UserInfo
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})

		return
	}
	// create user
	user := db.User{Username: body.Username, Email: body.Email, Password: string(hash)}
	result := db.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	sendToken(c, user)

	// send success response
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func Signin(c *gin.Context) {
	// get user data
	var body UserInfo
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// get user
	var user db.User
	result := db.DB.Where("email = ?", body.Email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	sendToken(c, user)
	// send response
	c.JSON(http.StatusOK, gin.H{
		"message": "User sign in successful",
	})
}

// Prevent users from accessing routes when unauthenicated.
func Protect(c *gin.Context) {
	token, err := c.Cookie("Authorization")
	if err == http.ErrNoCookie {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User is not logged in",
		})
		return
	}

	key := []byte(os.Getenv("JWT_SECRET"))

	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Could not validate auth token")
		}

		return key, nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User is not logged in",
		})
		return
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		// check token expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User is not logged in",
			})
			return
		}

		// get user
		var user db.User
		result := db.DB.Where("id = ?", claims["sub"]).First(&user)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User is not logged in",
			})
			return
		}

		// set user on req
		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "User is not logged in",
		})
		return
	}
}

// creates a jwt token for user
func sendToken(c *gin.Context, user db.User) {
	// creates a token struct with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
	})

	key := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	// send cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
}
