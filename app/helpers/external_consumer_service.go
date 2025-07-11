package helpers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// MakeRequest realiza una solicitud HTTP configurable y retorna el cuerpo y código de estado.
// method: método HTTP (GET, POST, etc.)
// requestUri: URL de destino
// body: cuerpo de la solicitud (puede ser nil)
// header: cabeceras adicionales
// isJson: indica si el cuerpo es JSON
// Retorna un mapa con el cuerpo de la respuesta y el código HTTP, o un error si ocurre

func MakeRequest(method string, requestUri string, body interface{}, header map[string]string, isJson bool) (map[string]interface{}, error) {
	// Preparar el cuerpo de la solicitud
	var requestBody []byte // Almacena el cuerpo serializado usando un slice de bytes
	var requestError error // Almacena errores de serialización

	if body != nil {
		// Si el cuerpo es JSON, serializarlo a JSON
		if isJson {
			requestBody, requestError = json.Marshal(body) // Serializa el cuerpo a JSON
		} else {
			// Si no es JSON, asumir que es un formulario y codificarlo
			formData := url.Values{} // Crear un nuevo objeto url.Values para almacenar los datos del formulario
			for key, value := range body.(map[string]string) {
				formData.Set(key, value)
			}
			requestBody = []byte(formData.Encode())
		}
		// Si hubo error al serializar, retornar el error
		if requestError != nil {
			return nil, requestError
		}
	}

	// Crear una nueva solicitud HTTP con el método, URL y cuerpo preparado
	httpRequest, httpRequestError := http.NewRequest(method, requestUri, bytes.NewBuffer(requestBody))
	if httpRequestError != nil {
		return nil, httpRequestError
	}

	// Configurar las cabeceras de la solicitud
	for key, value := range header {
		httpRequest.Header.Set(key, value)
	}

	// Si el cuerpo es JSON, establecer el tipo de contenido como application/json
	if isJson {
		httpRequest.Header.Set("Content-Type", "application/json")
	}

	// Crear un cliente HTTP y ejecutar la solicitud
	var client = &http.Client{}
	httpResponse, httpRequestError := client.Do(httpRequest)
	if httpRequestError != nil {
		return nil, httpRequestError
	}

	// Cerrar el cuerpo de la respuesta al finalizar, registrando error si ocurre
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(httpResponse.Body)

	// Leer el cuerpo de la respuesta
	bodyBytes, httpRequestError := io.ReadAll(httpResponse.Body)
	if httpRequestError != nil {
		return nil, httpRequestError
	}

	// Crear un mapa para incluir el cuerpo y el código HTTP en la respuesta
	responseMap := map[string]interface{}{
		"body":      string(bodyBytes),
		"http_code": httpResponse.StatusCode,
	}

	// Retornar el mapa con la respuesta y nil como error
	return responseMap, nil
}
