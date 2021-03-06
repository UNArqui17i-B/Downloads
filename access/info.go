package access

import (
	"net/http"
	"bytes"
	"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"os"
	"fmt"
)

func GetInformation(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileID := vars["id"]

	if fileID != "" {
		var buffer bytes.Buffer

		// Get complete document URL
		buffer.WriteString(os.Getenv("DB_URL") + "/" + os.Getenv("DB_NAME") + "/")
		buffer.WriteString(fileID)
		url := buffer.String()

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal("Connection to DB response: ", err)
			return
		}

		doc := new(FileInformation)
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}

		js, err := json.Marshal(doc)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		if resp.StatusCode == http.StatusOK {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(js)
			fmt.Println("Info request: 200")
		} else {
			rw.WriteHeader(http.StatusNotFound)
			fmt.Println("Info request: 404")
		}
	}
}