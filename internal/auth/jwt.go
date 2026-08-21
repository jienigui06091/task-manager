package auth
import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateJwtToken(username string,userId int64)(string,error){
	secret:= "qwueiwdhjdgudcbhuif"
claims := jwt.MapClaims{
		"sub":      userId,
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}