package proxy

import (
	"fmt"
	"io"
	"net/http"
)

// writeRequestHeaders writes the headers from request to writer.
func writeRequestHeaders(writer io.Writer, request *http.Request) error {
	if _, err := fmt.Fprintf(
		writer,
		"%s %s %s\r\n",
		request.Method,
		request.URL.RequestURI(),
		request.Proto,
	); err != nil {
		return err
	}

	if err := request.Header.Write(writer); err != nil {
		return err
	}

	_, err := io.WriteString(writer, "\r\n")
	return err
}

// writeResponseHeaders writes the headers from response to writer.
func writeResponseHeaders(writer http.ResponseWriter, response *http.Response) {
	headers := writer.Header()
	for name, values := range response.Header {
		if !isHopByHopHeader(name) {
			headers[name] = values
		}
	}

	if isWebSocketUpgrade(response.Header) {
		headers.Set("Connection", "upgrade")
		headers.Set("Upgrade", "websocket")
	}

	writer.WriteHeader(response.StatusCode)
}

// writeResponse writes the entirety of response to writer.
func writeResponse(writer http.ResponseWriter, response *http.Response) (int64, error) {
	defer response.Body.Close()
	writeResponseHeaders(writer, response)
	return io.Copy(writer, response.Body)
}
