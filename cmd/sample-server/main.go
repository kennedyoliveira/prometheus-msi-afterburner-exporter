// this is a sample mock server for the MSI Afterburner that will return a fixed set of metrics
// should be used for development or testing when you don't have any afterburner running at moment
// uses the default auth information, username=MSIAfterburner, password=17cc95b4017d496f82
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	listenAddress = flag.String("listen-address", "0.0.0.0:1082", "Host and port to listen to")
)

func main() {
	http.HandleFunc("/mahm", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("New request from %s %s %s ", request.RemoteAddr, request.Method, request.RequestURI)
		if request.Method != http.MethodGet {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// read from the responses used for test
		file, err := os.OpenFile("afterburner/test_samples/api_response.xml", os.O_RDONLY, 0)
		if err != nil {
			writer.WriteHeader(500)
			return
		}

		writer.Header().Add("Content-Type", "application/xml")

		_, copyErr := io.Copy(writer, file)
		if copyErr != nil {
			log.Printf("Failed to send response %v", copyErr)
			return
		}
	})

	log.Printf("Listening on: %s", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
