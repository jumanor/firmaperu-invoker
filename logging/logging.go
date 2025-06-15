package logging

import (
	"firmaperuweb/config"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log zerolog.Logger

func init() {

	abs_fname, _ := filepath.Abs("./")
	ruta := abs_fname + string(filepath.Separator) + "config.properties"

	properties, err := config.ReadPropertiesFile(ruta)
	if err != nil {
		panic(err)
	}

	var maxSize, maxBackups, maxAge int

	if properties["maxSize"] != "" {
		maxSize, err = strconv.Atoi(properties["maxSize"])
		if err != nil {
			panic(fmt.Sprintf("invalid maxSize value: %v", err))
		}
	}
	if properties["maxBackups"] != "" {
		maxSize, err = strconv.Atoi(properties["maxBackups"])
		if err != nil {
			panic(fmt.Sprintf("invalid maxBackups value: %v", err))
		}
	}
	if properties["maxAge"] != "" {
		maxSize, err = strconv.Atoi(properties["maxAge"])
		if err != nil {
			panic(fmt.Sprintf("invalid maxAge value: %v", err))
		}
	}

	err = os.MkdirAll("logs", 0755)
	if err != nil {
		panic(err)
	}

	logWriter := &lumberjack.Logger{
		Filename:   "logs/main.log", // Nombre del archivo de log
		MaxSize:    maxSize,         // Tamaño máximo en MB antes de rotar
		MaxBackups: maxBackups,      // Número máximo de archivos de backup
		MaxAge:     maxAge,          // Máximo de días a mantener los logs
		Compress:   true,            // Comprimir los archivos antiguos
	}

	var writers []io.Writer
	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	writers = append(writers, logWriter) //usamos el buffer del sistema operativo
	mw := io.MultiWriter(writers...)

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + fmt.Sprintf("%d", line)
	}

	log = zerolog.New(mw).With().Caller().Timestamp().Logger()
}
func Log() *zerolog.Logger {
	return &log
}
