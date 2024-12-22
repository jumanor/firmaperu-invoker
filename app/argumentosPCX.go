package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"firmaperuweb/logging"
	"firmaperuweb/util"

	"github.com/google/uuid"
)

type Pdf []struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// Definir la estructura para "firma"
type Firma struct {
	PosX            int    `json:"posx"`
	PosY            int    `json:"posy"`
	Reason          string `json:"reason"`
	Role            string `json:"role"`
	StampSigned     string `json:"stampSigned"`
	PageNumber      int    `json:"pageNumber"`
	VisiblePosition bool   `json:"visiblePosition"`
	SignatureStyle  int    `json:"signatureStyle"`
}

// Definir la estructura para "DatoArgumentos"
type DatoArgumentos struct {
	Pdfs  Pdf   `json:"pdfs"`
	Firma Firma `json:"firma"`
}

func ArgumentsServletPCX(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("x-access-token")
	if err := util.VerificarJWT(token); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Acceso no autorizado"))
		return
	}

	var inputParameter DatoArgumentos
	err := json.NewDecoder(r.Body).Decode(&inputParameter)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound) //codigo http 404
		w.Write([]byte("No se pudo parsear a json los parametros de entrada"))
		return
	}

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

	documentNameUUID, err := createFile7z(inputParameter.Pdfs)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound) //codigo http 404
		w.Write([]byte(err.Error()))
		return
	}

	param_query := "posx=" + strconv.Itoa(inputParameter.Firma.PosX) +
		"&posy=" + strconv.Itoa(inputParameter.Firma.PosY) +
		"&documentNameUUID=" + documentNameUUID +
		"&reason=" + url.QueryEscape(inputParameter.Firma.Reason) +
		"&role=" + url.QueryEscape(inputParameter.Firma.Role) +
		"&imageToStamp=" + url.QueryEscape(inputParameter.Firma.StampSigned) +
		"&visiblePosition=" + strconv.FormatBool(inputParameter.Firma.VisiblePosition) + //por defecto VisiblePosition=false
		"&signatureStyle=" + strconv.Itoa(inputParameter.Firma.SignatureStyle) + //por defecto SignatureStyle=0
		"&stampPage=" + strconv.Itoa(inputParameter.Firma.PageNumber) //por defecto PageNumber=0

	objetoJSON := map[string]string{
		"param_url":          serverURL + "/argumentos?" + param_query,
		"param_token":        token,
		"document_extension": "pdf",
	}

	// Convertir el mapa a una cadena JSON
	jsonBytes, err := json.Marshal(objetoJSON)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		return
	}

	// Convertir la cadena JSON a Base64
	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)

	urlBasePDFDownloadSigned := serverURL + "/downloadPdfSigned/" + url.QueryEscape(documentNameUUID)

	previewRespuesta := map[string]interface{}{
		"codigo": 2000, //codigo interno
		"data": map[string]interface{}{
			"argumentosBase64": base64Str,
			"urlBase":          urlBasePDFDownloadSigned,
		},
	}

	respuesta, _ := json.Marshal(previewRespuesta)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //codigo http 200
	w.Write(respuesta)

}

type ResultCanalDescarga struct {
	Message string
	Error   error
}

// Descargamos y creamos el documento PDF en el HD mediante una gorutina
func downloadPdfAndPersist(rutaMain string, pdf struct {
	URL  string "json:\"url\""
	Name string "json:\"name\""
}, ch chan ResultCanalDescarga) {

	out, err := os.Create(filepath.Join(rutaMain, pdf.Name+".pdf"))
	if err != nil {
		fmt.Println(err)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo crear archivo: " + pdf.Name)}
		return
	}

	client := http.Client{
		Timeout: 60 * time.Second, //timeout 60 segundos
	}
	resp, err := client.Get(pdf.URL)
	if err != nil {
		fmt.Println(err)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo descargar: " + pdf.URL)}
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("Error : " + resp.Status + "  " + pdf.URL)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo descargar: " + pdf.URL)}
		return

	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		ch <- ResultCanalDescarga{Message: "", Error: errors.New("No se pudo copiar a disco: " + pdf.URL)}
		return
	}

	out.Close()
	resp.Body.Close()

	//escritura en disco puede ser ineficiente
	//acceso concurrente a un unico archivo de escritura
	//el buffer del SO deberia de controlar este cuello de botella
	logging.Log().Debug().Str("nombre", pdf.Name).Str("url", pdf.URL).Msg("descargando")

	ch <- ResultCanalDescarga{Message: "OK " + pdf.Name, Error: nil}
}

// Descargamos todos los documentos PDF concurrentemente.
// Todos los documentos deben ser persistidos caso contrario se lanza en error
func downloadAllPdfAndPersistConcurrency(rutaMain string, urls Pdf) error {

	ch := make(chan ResultCanalDescarga)

	for _, pdf := range urls {
		go downloadPdfAndPersist(rutaMain, pdf, ch) //usamos go rutinas
	}

	for range urls {

		result := <-ch //bloqueamos a la espera de la respuesta

		if result.Error != nil {
			fmt.Println(result.Error)

			return result.Error
		}

		//fmt.Println(result.Message)

	}

	return nil
}

// Creamos un archivo 7z con los PDFs descargados
func createFile7z(urls Pdf) (string, error) {

	nameUUID := uuid.New().String()

	rutaMain := filepath.Join(os.TempDir(), "upload", nameUUID)

	if err := os.MkdirAll(rutaMain, os.ModePerm); err != nil {
		logging.Log().Error().Err(err).Send()
		return "", errors.New("No se puede crear el directorio " + rutaMain)
	}

	if err := downloadAllPdfAndPersistConcurrency(rutaMain, urls); err != nil {
		logging.Log().Error().Err(err).Send()
		return "", err
	}

	file7z := filepath.Join(rutaMain, "..", nameUUID+".7z")
	c := exec.Command("7z", "a", file7z, rutaMain+string(filepath.Separator)+".")

	if err := c.Run(); err != nil {
		logging.Log().Error().Err(err).Send()
		return "", errors.New("no se pudo comprimir a 7z")
	}

	logging.Log().Debug().Str("7z", file7z).Msg("Archivo 7z creado satisfactoriamente")
	return nameUUID, nil
}
