package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestLoginFailures(t *testing.T) {

	t.Run("LoginWithInvalidUser", func(t *testing.T) {
		loginReq := LoginRequest{
			UsuarioAccesoApi: "usuarioInvalido",
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

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Estado esperado 404, se recibió %d", resp.StatusCode)
		}

		// Captura y verifica el mensaje de error
		var errorMsg []byte
		errorMsg, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error al leer la respuesta: %v", err)
		}

		expectedError := "Usuario incorrecto"
		if string(errorMsg) != expectedError {
			t.Errorf("Mensaje de error esperado '%s', se recibió: '%s'", expectedError, string(errorMsg))
		}

		t.Logf("Respuesta correcta para usuario inválido: %d", resp.StatusCode)
	})

	t.Run("LoginWithEmptyUser", func(t *testing.T) {
		loginReq := LoginRequest{
			UsuarioAccesoApi: "",
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

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Estado esperado 404, se recibió %d", resp.StatusCode)
		}

		// Captura y verifica el mensaje de error
		var errorMsg []byte
		errorMsg, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error al leer la respuesta: %v", err)
		}

		expectedError := "Usuario incorrecto"
		if string(errorMsg) != expectedError {
			t.Errorf("Mensaje de error esperado '%s', se recibió: '%s'", expectedError, string(errorMsg))
		}

		t.Logf("Respuesta correcta para usuario vacío: %d", resp.StatusCode)
	})

	t.Run("LoginWithMalformedJSON", func(t *testing.T) {
		malformedJSON := `{"usuarioAccesoApi": "usuario"` // JSON incompleto

		resp, err := http.Post(
			baseURLFirmaPeruInvoker+"/autenticacion",
			"application/json",
			bytes.NewBuffer([]byte(malformedJSON)),
		)
		if err != nil {
			t.Fatalf("Error al realizar la solicitud de login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Estado esperado 404, se recibió %d", resp.StatusCode)
		}

		// Captura y verifica el mensaje de error
		var errorMsg []byte
		errorMsg, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error al leer la respuesta: %v", err)
		}

		expectedError := "No se pudo parsear a json los parametros de entrada"
		if string(errorMsg) != expectedError {
			t.Errorf("Mensaje de error esperado '%s', se recibió: '%s'", expectedError, string(errorMsg))
		}

		t.Logf("Respuesta correcta para JSON malformado: %d", resp.StatusCode)
	})

	t.Run("LoginWithMissingField", func(t *testing.T) {
		// JSON sin el campo requerido
		jsonData := `{}`

		resp, err := http.Post(
			baseURLFirmaPeruInvoker+"/autenticacion",
			"application/json",
			bytes.NewBuffer([]byte(jsonData)),
		)
		if err != nil {
			t.Fatalf("Error al realizar la solicitud de login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Estado esperado 404, se recibió %d", resp.StatusCode)
		}

		// Captura y verifica el mensaje de error
		var errorMsg []byte
		errorMsg, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error al leer la respuesta: %v", err)
		}

		expectedError := "Usuario incorrecto"
		if string(errorMsg) != expectedError {
			t.Errorf("Mensaje de error esperado '%s', se recibió: '%s'", expectedError, string(errorMsg))
		}

		t.Logf("Respuesta correcta para campo faltante: %d", resp.StatusCode)
	})

	t.Run("LoginWithNullUser", func(t *testing.T) {
		// JSON con valor null
		jsonData := `{"usuarioAccesoApi": null}`

		resp, err := http.Post(
			baseURLFirmaPeruInvoker+"/autenticacion",
			"application/json",
			bytes.NewBuffer([]byte(jsonData)),
		)
		if err != nil {
			t.Fatalf("Error al realizar la solicitud de login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Estado esperado 404, se recibió %d", resp.StatusCode)
		}

		// Captura y verifica el mensaje de error
		var errorMsg []byte
		errorMsg, err = io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Error al leer la respuesta: %v", err)
		}

		expectedError := "Usuario incorrecto"
		if string(errorMsg) != expectedError {
			t.Errorf("Mensaje de error esperado '%s', se recibió: '%s'", expectedError, string(errorMsg))
		}

		t.Logf("Respuesta correcta para usuario null: %d", resp.StatusCode)
	})
}
