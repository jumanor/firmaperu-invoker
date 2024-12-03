package app

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"firmaperuweb/logging"
	"firmaperuweb/util"
)

// Exclusivamente utilizado por ReFirmaPCX para descargar los documentos (sin firmar) que esta comprimidos con 7z
func Download7z(w http.ResponseWriter, r *http.Request) {

	token := r.URL.Query().Get("token")
	if err := util.VerificarJWT(token); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	documentName7z := r.URL.Query().Get("documentName") + ".7z"
	filename := filepath.Join(os.TempDir(), "upload", documentName7z)
	logging.Log().Debug().Str("7z", filename).Msg("Archivo 7z descargado satisfactoriamente")

	// Open file
	f, err := os.Open(filename)
	if err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//Set header
	w.Header().Add("Content-type", "application/octet-stream")
	w.Header().Add("Content-disposition", "attachment; filename="+documentName7z)

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		logging.Log().Error().Err(err).Send()
		w.WriteHeader(500)
	}
}
