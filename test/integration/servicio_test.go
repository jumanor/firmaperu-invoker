package integration_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"testing"
)

const baseURLFirmaPeruInvoker = "http://127.0.0.1:9091"
const baseURLServidorWebExamples = "http://127.0.0.1:80"

type LoginRequest struct {
	UsuarioAccesoApi string `json:"UsuarioAccesoApi"`
}

type LoginResponse struct {
	Data string `json:"data"`
}

type ArgumentosRequest struct {
	PDFs  []PDFInfo `json:"pdfs"`
	Firma Firma     `json:"firma"`
}

type PDFInfo struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

type Firma struct {
	PosX int `json:"posx"`
	PosY int `json:"posy"`
}

type ArgumentosResponse struct {
	Data struct {
		ArgumentosBase64 string `json:"argumentosBase64"`
		URLBase          string `json:"urlBase"`
	} `json:"data"`
}

type DownloadPDFBase64Response struct {
	Codigo int    `json:"codigo"`
	Data   string `json:"data"`
}

func TestIntegrationSequential(t *testing.T) {
	var authTokenGlobal string

	t.Run("Login", func(t *testing.T) {
		loginReq := LoginRequest{
			UsuarioAccesoApi: "usuarioAccesoApi",
		}

		jsonData, err := json.Marshal(loginReq)
		if err != nil {
			t.Fatalf("Error al serializar la solicitud de login: %v", err)
		}

		resp, err := http.Post(
			baseURLFirmaPeruInvoker+"/autenticacion",
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			t.Fatalf("Error al realizar la solicitud de login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado 200, se recibió %d", resp.StatusCode)
		}

		var loginResp LoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			t.Fatalf("Error al decodificar la respuesta de login: %v", err)
		}

		if loginResp.Data == "" {
			t.Error("Se esperaba un token de autenticación, se recibió una cadena vacía")
		}

		// Guardar el token para el siguiente test
		authTokenGlobal = loginResp.Data
		t.Logf("Token de autenticación: %s", loginResp.Data)

		// Verificar formato JWT
		tokenParts := strings.Split(loginResp.Data, ".")
		if len(tokenParts) != 3 {
			t.Errorf("Se esperaba formato JWT (cabecera.carga.firma), se obtuvieron %d partes", len(tokenParts))
		}
	})

	var argumentosBase64Global string
	var urlBaseGlobal string

	t.Run("ArgumentosServletPCX", func(t *testing.T) {

		authToken := authTokenGlobal

		argumentosReq := ArgumentosRequest{
			PDFs: []PDFInfo{
				{URL: baseURLServidorWebExamples + "/example01/01.pdf", Name: "doc0"},
				{URL: baseURLServidorWebExamples + "/example01/02.pdf", Name: "doc1"},
			},
			Firma: Firma{
				PosX: 10,
				PosY: 12,
			},
		}

		jsonData, err := json.Marshal(argumentosReq)
		if err != nil {
			t.Fatalf("Error al serializar la solicitud de argumentos: %v", err)
		}

		req, err := http.NewRequest("POST", baseURLFirmaPeruInvoker+"/argumentsServletPCX", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Error al crear la solicitud: %v", err)
		}

		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		req.Header.Set("x-access-token", authToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error al realizar la solicitud de argumentos: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado 200, se recibió %d", resp.StatusCode)
		}

		var argumentosResp ArgumentosResponse
		if err := json.NewDecoder(resp.Body).Decode(&argumentosResp); err != nil {
			t.Fatalf("Error al decodificar la respuesta de argumentos: %v", err)
		}

		if argumentosResp.Data.ArgumentosBase64 == "" {
			t.Error("Se esperaba argumentosBase64, se recibió una cadena vacía")
		}

		if argumentosResp.Data.URLBase == "" {
			t.Error("Se esperaba urlBase, se recibió una cadena vacía")
		}

		t.Logf("ArgumentosBase64: %s", argumentosResp.Data.ArgumentosBase64)
		t.Logf("URLBase: %s", argumentosResp.Data.URLBase)

		argumentosBase64Global = argumentosResp.Data.ArgumentosBase64 // Guardar para el siguiente test
		urlBaseGlobal = argumentosResp.Data.URLBase                   // Guardar para el siguiente test
	})

	var documentUUIDGlobal string

	t.Run("TestDownload7z", func(t *testing.T) {
		authToken := authTokenGlobal
		argumentosBase64 := argumentosBase64Global
		// Decode Base64 string to validate it's properly encoded
		decodedArguments, err := base64.StdEncoding.DecodeString(argumentosBase64)
		if err != nil {
			t.Fatalf("Error al decodificar argumentos en Base64: %v", err)
		}

		//map[document_extension:pdf param_token:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJqdW1hbm9yIiwiZXhwIjoxNzU5MzU4MzA5fQ.61nGOUzVP_J96Z6pXOLPLprMZTFfcZrMbRZ0eiinCYw param_url:http://127.0.0.1:9091/argumentos?posx=10&posy=12&documentNameUUID=bbfed18b-c025-4911-bb8e-6077f047593f&reason=&role=&imageToStamp=&visiblePosition=false&oneByOne=false&signatureStyle=-1&stampPage=0&stampTextSize=0&stampWordWrap=0]
		var parsedArgs map[string]interface{}
		if err := json.Unmarshal(decodedArguments, &parsedArgs); err != nil {
			t.Fatalf("Error al parsear argumentos decodificados: %v", err)
		}

		paramUrl, ok := parsedArgs["param_url"].(string)
		if !ok || paramUrl == "" {
			t.Fatal("No se encontró param_token o es inválido en los argumentos parseados")
		}

		//http://127.0.0.1:9091/argumentos?posx=10&posy=12&documentNameUUID=d3f825a2-81e3-4429-b4a5-738d44774a75&reason=&role=&imageToStamp=&visiblePosition=false&oneByOne=false&signatureStyle=-1&stampPage=0&stampTextSize=0&stampWordWrap=
		t.Logf("param_url: %s", paramUrl)

		urlParts := strings.Split(paramUrl, "?")
		if len(urlParts) != 2 {
			t.Fatal("Formato de param_url inválido")
		}
		queryParams := urlParts[1]
		t.Logf("Parámetros de consulta: %s", queryParams)

		var documentUUID string
		for _, param := range strings.Split(queryParams, "&") {
			kv := strings.SplitN(param, "=", 2)
			if len(kv) == 2 && kv[0] == "documentNameUUID" {
				documentUUID = kv[1]
				break
			}
		}

		if documentUUID == "" {
			t.Fatal("No se encontró documentNameUUID en los parámetros de consulta")
		}

		documentUUIDGlobal = documentUUID

		t.Logf("documentNameUUID extraído: %s", documentUUID)

		url := fmt.Sprintf("%s/download7z?documentName=%s&token=%s", baseURLFirmaPeruInvoker, documentUUID, authToken)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("Error al crear solicitud de descarga: %v", err)
		}

		req.Header.Set("Accept", "application/octet-stream")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error al realizar solicitud de descarga: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado 200, se recibió %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/octet-stream") {
			t.Errorf("Se esperaba tipo de contenido octet-stream, se recibió %s", contentType)
		}

		// Verificar que hay contenido
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error al leer el cuerpo de la respuesta: %v", err)
		}

		if len(body) == 0 {
			t.Error("Se esperaba contenido del archivo, se recibió respuesta vacía")
		}

		t.Logf("Tamaño del archivo descargado: %d bytes", len(body))
	})

	t.Run("TestUpload7z", func(t *testing.T) {
		documentUUID := documentUUIDGlobal
		authToken := authTokenGlobal

		// Abrir el archivo 7z real
		file, err := os.Open("../api/firmados.7z")
		if err != nil {
			t.Fatalf("Error al abrir el archivo firmados.7z: %v", err)
		}
		defer file.Close()

		// Crear buffer y writer para multipart form
		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)

		// Crear parte del formulario para el archivo
		part, err := writer.CreateFormFile("archivo7z", "firmados.7z")
		if err != nil {
			t.Fatalf("Error al crear el campo de archivo en el formulario: %v", err)
		}

		// Copiar contenido del archivo al form
		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("Error al copiar el contenido del archivo: %v", err)
		}

		writer.Close()

		url := fmt.Sprintf("%s/upload7z/%s?token=%s", baseURLFirmaPeruInvoker, documentUUID, authToken)

		req, err := http.NewRequest("POST", url, &buf)
		if err != nil {
			t.Fatalf("Error al crear la solicitud de carga: %v", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Error al realizar la solicitud de carga: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado 200, se recibió %d", resp.StatusCode)
		}

		t.Logf("Carga completada con estado: %d", resp.StatusCode)
	})

	t.Run("TestDownloadSignedPDF", func(t *testing.T) {

		authToken := authTokenGlobal
		namePdf := "doc0"
		urlBase := urlBaseGlobal

		// Prueba método GET
		t.Run("Método GET", func(t *testing.T) {
			url := fmt.Sprintf("%s/%s/%s", urlBase, namePdf, authToken)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("Error al crear la solicitud GET: %v", err)
			}

			req.Header.Set("Accept", "application/pdf")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error al realizar la solicitud GET: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Estado esperado 200, se recibió %d", resp.StatusCode)
			}

			// Leer todo el contenido de la respuesta
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error al leer el cuerpo de la respuesta: %v", err)
			}

			if len(body) == 0 {
				t.Error("Se esperaba contenido del archivo PDF, se recibió respuesta vacía")
			}

			t.Logf("Tamaño del PDF descargado: %d bytes", len(body))
		})

		// Prueba método POST
		t.Run("Método POST", func(t *testing.T) {
			url := fmt.Sprintf("%s/%s", urlBase, namePdf)

			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				t.Fatalf("Error al crear la solicitud POST: %v", err)
			}

			req.Header.Set("Accept", "application/pdf")
			req.Header.Set("x-access-token", authToken)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error al realizar la solicitud POST: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Estado esperado 200, se recibió %d", resp.StatusCode)
			}

			// Verificar Content-Type para JSON
			contentType := resp.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				t.Errorf("Se esperaba Content-Type application/json, se recibió %s", contentType)
			}

			// Leer y parsear la respuesta JSON
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Error al leer el cuerpo de la respuesta: %v", err)
			}

			if len(body) == 0 {
				t.Error("Se esperaba contenido JSON, se recibió respuesta vacía")
			}

			var pdfResp DownloadPDFBase64Response
			if err := json.Unmarshal(body, &pdfResp); err != nil {
				t.Fatalf("Error al parsear respuesta JSON: %v", err)
			}

			// Verificar estructura de la respuesta
			if pdfResp.Codigo != 2000 {
				t.Errorf("Se esperaba código 2000, se recibió %d", pdfResp.Codigo)
			}

			if pdfResp.Data == "" {
				t.Error("Se esperaba data con contenido Base64, se recibió cadena vacía")
			}

			// Verificar que el data es válido Base64
			decodedPDF, err := base64.StdEncoding.DecodeString(pdfResp.Data)
			if err != nil {
				t.Fatalf("Error al decodificar Base64: %v", err)
			}

			if len(decodedPDF) == 0 {
				t.Error("Se esperaba contenido PDF decodificado, se recibió contenido vacío")
			}

			// Verificar que parece ser un PDF (debe comenzar con %PDF)
			if len(decodedPDF) >= 4 {
				pdfHeader := string(decodedPDF[:4])
				if pdfHeader != "%PDF" {
					t.Errorf("Se esperaba header PDF (%%PDF), se recibió %s", pdfHeader)
				}
			}

			t.Logf("Respuesta JSON - Código: %d", pdfResp.Codigo)
			t.Logf("Tamaño del PDF en Base64: %d caracteres", len(pdfResp.Data))
			t.Logf("Tamaño del PDF decodificado: %d bytes", len(decodedPDF))
		})
	})
}
