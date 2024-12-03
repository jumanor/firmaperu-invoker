package app

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"firmaperuweb/logging"
	"firmaperuweb/util"

	"github.com/gorilla/mux"
)

// Descargamos el documento PDF firmado mediante GET
func DownloadPdfSigned(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	namePdf, err := url.QueryUnescape(vars["file"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de " + vars["file"]))
		return
	}
	dirPdf, err := url.QueryUnescape(vars["dir"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de  " + vars["dir"]))
		return
	}

	token, err := url.QueryUnescape(vars["token"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de  " + vars["token"]))
		return
	}
	if err := util.VerificarJWT(token); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Acceso no autorizado"))
		return
	}

	filePdfSigned := filepath.Join(os.TempDir(), "upload", "signed", dirPdf+"[R]", namePdf+"[FP].pdf")

	// Open file
	f, err := os.Open(filePdfSigned)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo abrir el archivo" + namePdf))
		return
	}
	defer f.Close()

	//Set header
	w.Header().Add("Content-type", "application/pdf")
	w.Header().Add("Content-disposition", "filename="+namePdf+"[R].pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("No se pudo copiar el archivo " + namePdf + " en flujo de envio "))
	} else {
		logging.Log().Debug().Str("nombre", namePdf).Str("ruta", filePdfSigned).Msg("descargando")
	}

}

// Descargamos el docuemento PDF firmado en Base64 mediante POST
func DownloadPdfSignedBase64(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	namePdf, err := url.QueryUnescape(vars["file"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de " + vars["file"]))
		return
	}
	dirPdf, err := url.QueryUnescape(vars["dir"])
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo realizar decode de  " + vars["dir"]))
		return
	}

	token := r.Header.Get("x-access-token")

	if err := util.VerificarJWT(token); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Acceso no autorizado"))
		return
	}

	filePdfSigned := filepath.Join(os.TempDir(), "upload", "signed", dirPdf+"[R]", namePdf+"[FP].pdf")

	data, err := os.ReadFile(filePdfSigned)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No se pudo leer el archivo" + namePdf))
		return
	}

	logging.Log().Debug().Str("nombre", namePdf).Str("ruta", filePdfSigned).Msg("descargando b64")

	previewRespuesta := map[string]interface{}{
		"codigo": 2000, //codigo interno
		"data":   base64.StdEncoding.EncodeToString(data),
	}

	respuesta, _ := json.Marshal(previewRespuesta)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) //codigo http 200
	w.Write(respuesta)

}
