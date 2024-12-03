package util

import (
	"firmaperuweb/logging"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var SECRET_KEY_JWT string
var TIME_EXPIRE_TOKEN int64

// Creamos Token
func GenerarJWT() (string, error) {

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(TIME_EXPIRE_TOKEN) * time.Minute)),
		Issuer:    "jumanor",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(SECRET_KEY_JWT))
}

// Verificamos Token
func VerificarJWT(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY_JWT), nil
	})

	if token.Valid {
		return nil
	} else {
		return err
	}

}

func VerificarExpiracionJWT(tokenString string) bool {

	// Decodificar el token sin verificar la firma
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		logging.Log().Error().Err(err).Msg("Error al decodificar el token")
		return true
	}

	// Obtener las reclamaciones (claims)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Extraer la fecha de expiración (exp)
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			//fmt.Println("Fecha de caducidad:", expirationTime)
			//fmt.Println("Expirado:", time.Now().After(expirationTime))
			logging.Log().Debug().Msg("El token aun no expira")
			return time.Now().After(expirationTime)

		} else {

			logging.Log().Error().Err(err).Msg("El token no tiene un campo **exp**")
			return true
		}

		// Mostrar todos los claims (opcional)
		//claimsJSON, _ := json.MarshalIndent(claims, "", "  ")
		//fmt.Println("Claims:", string(claimsJSON))
	} else {
		logging.Log().Error().Err(err).Msg("No se pudieron obtener las claims del token")
		return true
	}
}
