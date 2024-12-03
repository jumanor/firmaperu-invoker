package main

import (
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"firmaperuweb/app"
	"firmaperuweb/logging"
	"firmaperuweb/util"

	"github.com/gorilla/mux"
)

var SERVER_ADDRESS string
var CERTIFICATE_FILE_TLS string
var PRIVATE_KEY_FILE_TLS string

func init() {

	util.TIME_EXPIRE_TOKEN = 5

	abs_fname, _ := filepath.Abs("./")
	ruta := abs_fname + string(filepath.Separator) + "config.properties"

	properties, err := util.ReadPropertiesFile(ruta)
	if err != nil {
		panic(err)
	}

	app.CLIENT_ID = properties["clientId"]
	app.CLIENT_SECRET = properties["clientSecret"]
	SERVER_ADDRESS = properties["serverAddress"]
	util.SECRET_KEY_JWT = properties["secretKeyJwt"]
	app.USER_ACCESS_API = properties["userAccessApi"]

	if properties["timeExpireToken"] != "" {
		if exp, err := strconv.ParseInt(properties["timeExpireToken"], 10, 64); err != nil {
			panic(err)
		} else {
			util.TIME_EXPIRE_TOKEN = exp
		}
	}
	//if properties["maxFileSize7z"] != "" {
	//	app.MAX_FILE_SIZE_7Z = properties["maxFileSize7z"]
	//}
	if properties["certificateFileTls"] != "" {
		CERTIFICATE_FILE_TLS = properties["certificateFileTls"]
	}
	if properties["privateKeyFileTls"] != "" {
		PRIVATE_KEY_FILE_TLS = properties["privateKeyFileTls"]
	}
}

func main() {
	enrutador := mux.NewRouter()

	fs := http.FileServer(http.Dir("./public/"))
	enrutador.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs))

	enrutador.HandleFunc("/argumentos", app.Argumentos).Methods("POST")
	enrutador.HandleFunc("/argumentsServletPCX", app.ArgumentsServletPCX).Methods("POST")
	enrutador.HandleFunc("/download7z", app.Download7z).Methods("GET")
	enrutador.HandleFunc("/upload7z/{uuid}", app.Upload7z).Methods("POST")
	enrutador.HandleFunc("/downloadPdfSigned/{dir}/{file}", app.DownloadPdfSignedBase64).Methods("POST")
	enrutador.HandleFunc("/downloadPdfSigned/{dir}/{file}/{token}", app.DownloadPdfSigned).Methods("GET")
	enrutador.HandleFunc("/autenticacion", app.Autenticacion).Methods("POST")
	//
	enrutador.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)

	enrutador.Use(util.EnableCors)

	servidor := &http.Server{
		Handler: enrutador,
		Addr:    SERVER_ADDRESS,
		// Timeouts para evitar que el servidor se quede "colgado" por siempre (2 minutos)
		WriteTimeout: 120 * time.Second,
		ReadTimeout:  120 * time.Second,
	}

	if CERTIFICATE_FILE_TLS != "" && PRIVATE_KEY_FILE_TLS != "" {

		logging.Log().Info().Str("Scheme", "https").Msgf("Escuchando en %s. Presiona CTRL + C para salir", SERVER_ADDRESS)
		err := servidor.ListenAndServeTLS(CERTIFICATE_FILE_TLS, PRIVATE_KEY_FILE_TLS)
		logging.Log().Fatal().Err(err).Send()

	} else {

		logging.Log().Info().Str("Scheme", "http").Msgf("Escuchando en %s. Presiona CTRL + C para salir", SERVER_ADDRESS)
		err := servidor.ListenAndServe()
		logging.Log().Fatal().Err(err).Send()
	}

}
