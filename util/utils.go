package util

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func ConfigCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,x-access-token")
}
func EnableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ConfigCors(&w)
		next.ServeHTTP(w, r)

	})
}

func CreateVersionFile(version, buildTime, gitCommit string) error {
	content := fmt.Sprintf("Version: %s\nBuildTime: %s\nGitCommit: %s\n", version, buildTime, gitCommit)
	return os.WriteFile(filepath.Join(".", "version.txt"), []byte(content), 0644)
}
