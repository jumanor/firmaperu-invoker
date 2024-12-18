package app

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"firmaperuweb/logging"
	"firmaperuweb/util"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var CLIENT_ID string
var CLIENT_SECRET string

// Exclusivamente utilizado por FirmaPeru
func Argumentos(w http.ResponseWriter, r *http.Request) {
	logging.Log().Trace().Msg("Inicio solicitando cadena de argumentos base64")

	scheme := "http"
	if proto := r.Header.Get("X-Forwarded-Proto"); proto == "https" {
		logging.Log().Debug().Str("request", "proxy ssl").Send()
		scheme = "https"
	}

	if r.TLS != nil {
		logging.Log().Debug().Str("request", "directo ssl").Send()
		scheme = "https"
	}
	serverURL := scheme + "://" + r.Host

	logging.Log().Debug().Str("url", serverURL).Msg("Ruta de construccion")

	param_token := r.FormValue("param_token")
	// Obtenemos parametros url query
	documentNameUUID := r.URL.Query().Get("documentNameUUID")
	posx := r.URL.Query().Get("posx")
	posy := r.URL.Query().Get("posy")
	signatureReason, _ := url.QueryUnescape(r.URL.Query().Get("reason"))
	role, _ := url.QueryUnescape(r.URL.Query().Get("role"))
	imageToStamp, _ := url.QueryUnescape(r.URL.Query().Get("imageToStamp"))
	stampPageQuery, _ := url.QueryUnescape(r.URL.Query().Get("stampPage"))
	visiblePositionQuery, _ := url.QueryUnescape(r.URL.Query().Get("visiblePosition"))
	signatureStyleQuery, _ := url.QueryUnescape(r.URL.Query().Get("signatureStyle"))
	//
	if imageToStamp == "" {
		imageToStamp = serverURL + "/public/iFirma.png"
	}
	stampPage := 1
	if stampPageQuery != "0" {
		var err error
		stampPage, err = strconv.Atoi(stampPageQuery)
		if err != nil {
			msn := "Error al convertir a entero variable stampPage (pageNumber)"
			logging.Log().Error().Err(err).Msg(msn)
			http.Error(w, msn, http.StatusInternalServerError) //codigo http 500
			return
		}
	}
	visiblePosition := false
	if visiblePositionQuery != "false" {
		var err error
		visiblePosition, err = strconv.ParseBool(visiblePositionQuery)
		if err != nil {
			msn := "Error al convertir a bool variable visiblePosition"
			logging.Log().Error().Err(err).Msg(msn)
			http.Error(w, msn, http.StatusInternalServerError) //codigo http 500
			return
		}
	}
	signatureStyle := 1
	if signatureStyleQuery != "0" {
		var err error
		signatureStyle, err = strconv.Atoi(signatureStyleQuery)
		if err != nil {
			msn := "Error al convertir a entero variable signatureStyle"
			logging.Log().Error().Err(err).Msg(msn)
			http.Error(w, msn, http.StatusInternalServerError) //codigo http 500
			return
		}
	}
	//
	documentToSign := serverURL + "/download7z?documentName=" + url.QueryEscape(documentNameUUID) + "&token=" + param_token
	uploadDocumentSigned := serverURL + "/upload7z/" + url.QueryEscape(documentNameUUID) + "?token=" + param_token
	//

	token_firma_peru, err := obtenerTokenFirmaPeru()
	if err != nil {
		msn := "Error al recuperar el Token Firma Peru"
		logging.Log().Error().Err(err).Msg(msn)
		http.Error(w, msn, http.StatusInternalServerError) //codigo http 500
		return
	}

	base64Str, err := convertir(signatureStyle, visiblePosition, stampPage, role, signatureReason, imageToStamp, documentToSign, uploadDocumentSigned, posx, posy, token_firma_peru)
	if err != nil {
		msn := "Error al convertir argumentos en base64"
		logging.Log().Error().Err(err).Msg(msn)
		http.Error(w, msn, http.StatusInternalServerError) //codigo http 500
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK) //codigo http 200
	w.Write([]byte(base64Str))
}

func readTokenFromFile() (string, error) {
	file, err := os.Open("token.txt")
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("error al leer el archivo: %v", err)
	}
	return string(content), nil
}

func saveTokenToFile(token string) error {
	file, err := os.Create("token.txt")
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}
	return nil
}
func obtenerTokenFirmaPeru() (string, error) {

	token, err := readTokenFromFile()
	if err != nil { //el token no existe en archivo

		token, err = crearTokenFirmaPeru()
		if err != nil {

			return "", err
		}
		err = saveTokenToFile(token)
		if err != nil {
			return "", err
		}
		return token, nil

	} else { //el token existe en archivo

		estado := util.VerificarExpiracionJWT(token)
		if !estado {

			logging.Log().Debug().Msg("Token Firma Peru aun no expira")
			return token, nil
		} else {
			//el token expiro
			token, err = crearTokenFirmaPeru()
			if err != nil {
				return "", err
			}
			err = saveTokenToFile(token)
			if err != nil {
				return "", err
			}
			return token, nil
		}
	}
}

// genera token firma peru
func crearTokenFirmaPeru() (string, error) {

	apiURL := "https://apps.firmaperu.gob.pe/admin/api/security/generate-token"
	formData := url.Values{
		"client_id":     {CLIENT_ID},
		"client_secret": {CLIENT_SECRET},
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {

		logging.Log().Error().Err(err).Msg("Error al crear la solicitud")
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		logging.Log().Error().Err(err).Msg("Error al enviar la solicitud")
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		logging.Log().Error().Err(err).Msg("Error al leer la respuesta")
		return "", err
	}

	logging.Log().Debug().Msg("Se genero un nuevo Token Firma Peru")
	return string(body), nil

}

func convertir(signatureStyle int, visiblePosition bool, stampPage int, role string, signatureReason string,
	imageToStamp string, documentToSign string,
	uploadDocumentSigned string, posx string, posy string, token string) (string, error) {

	param := map[string]interface{}{
		"signatureFormat":        "PAdES",
		"signatureLevel":         "B",
		"signaturePackaging":     "enveloped",
		"documentToSign":         documentToSign,
		"certificateFilter":      ".*",
		"webTsa":                 "",
		"userTsa":                "",
		"passwordTsa":            "",
		"theme":                  "claro",
		"visiblePosition":        visiblePosition, //importante
		"contactInfo":            "",
		"signatureReason":        signatureReason,
		"bachtOperation":         true,
		"oneByOne":               false,
		"signatureStyle":         signatureStyle, //default 1
		"imageToStamp":           imageToStamp,
		"stampTextSize":          14,
		"stampWordWrap":          37,
		"role":                   role,
		"stampPage":              stampPage, //default 1
		"positionx":              posx,
		"positiony":              posy,
		"uploadDocumentSigned":   uploadDocumentSigned,
		"certificationSignature": false, //unico firmante
		"token":                  token,
	}

	jsonBytes, err := json.Marshal(param)
	if err != nil {

		logging.Log().Fatal().Err(err).Msg("Error al convertir el JSON")
		return "", err
	}

	// Codificar en Base64
	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)

	return base64Str, nil

}
