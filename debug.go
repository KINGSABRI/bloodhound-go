package bloodhound

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/pretty"
)

var (
	// Style for pretty-printing JSON
	prettyStyle = &pretty.Style{
		Key:    [2]string{"\x1b[36m", "\x1b[0m"}, // Cyan
		String: [2]string{"\x1b[32m", "\x1b[0m"}, // Green
		Number: [2]string{"\x1b[33m", "\x1b[0m"}, // Yellow
		True:   [2]string{"\x1b[35m", "\x1b[0m"}, // Magenta
		False:  [2]string{"\x1b[35m", "\x1b[0m"}, // Magenta
		Null:   [2]string{"\x1b[31m", "\x1b[0m"}, // Red
	}
)

// logRequest is a helper to elegantly print HTTP request details.
func (c *Client) logRequest(req *http.Request, body []byte) {
	if c.debugWriter == nil {
		return
	}

	var out strings.Builder

	// Request Line
	out.WriteString(fmt.Sprintf("\x1b[1m--> %s %s %s\x1b[0m\n", req.Method, req.URL.RequestURI(), req.Proto))

	// Headers
	for name, headers := range req.Header {
		for _, h := range headers {
			out.WriteString(fmt.Sprintf("\x1b[36m%s\x1b[0m: %s\n", name, h))
		}
	}
	out.WriteString("\n")

	// Body
	if len(body) > 0 {
		formattedBody := pretty.Color(pretty.Pretty(body), prettyStyle)
		out.Write(formattedBody)
		out.WriteString("\n")
	}

	out.WriteString("\x1b[1m--> END %s\x1b[0m\n\n")
	fmt.Fprint(c.debugWriter, out.String())
}

// logResponse is a helper to elegantly print HTTP response details.
func (c *Client) logResponse(resp *http.Response) {
	if c.debugWriter == nil {
		return
	}

	var out strings.Builder

	// Status Line
	statusColor := "\x1b[32m" // Green for 2xx
	if resp.StatusCode >= 400 {
		statusColor = "\x1b[31m" // Red for 4xx/5xx
	} else if resp.StatusCode >= 300 {
		statusColor = "\x1b[33m" // Yellow for 3xx
	}
	out.WriteString(fmt.Sprintf("\x1b[1m<-- %s %s\x1b[0m\n", resp.Proto, statusColor+resp.Status+"\x1b[0m"))

	// Headers
	for name, headers := range resp.Header {
		for _, h := range headers {
			out.WriteString(fmt.Sprintf("\x1b[36m%s\x1b[0m: %s\n", name, h))
		}
	}
	out.WriteString("\n")

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		out.WriteString(fmt.Sprintf("\x1b[31mError reading response body: %v\x1b[0m\n", err))
	} else {
		// Decompress if gzipped
		var uncompressedBody []byte
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(bytes.NewReader(body))
			if err != nil {
				out.WriteString(fmt.Sprintf("\x1b[31mError creating gzip reader: %v\x1b[0m\n", err))
			} else {
				uncompressedBody, _ = io.ReadAll(reader)
				reader.Close()
			}
		} else {
			uncompressedBody = body
		}

		// Pretty-print the JSON body
		if len(uncompressedBody) > 0 {
			formattedBody := pretty.Color(pretty.Pretty(uncompressedBody), prettyStyle)
			out.Write(formattedBody)
			out.WriteString("\n")
		}
		// Restore the original (potentially compressed) body for the main logic to re-read and decompress.
		resp.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	out.WriteString("\x1b[1m<-- END HTTP\x1b[0m\n\n")
	fmt.Fprint(c.debugWriter, out.String())
}
