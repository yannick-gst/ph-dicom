package main

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

	"github.com/google/uuid"
    "github.com/suyashkumar/dicom"
)

var destinationDir string

func createDestinationDir() (string, error) {
    _, err := os.Stat("./dicom/out")
    if os.IsNotExist(err) {
        // The directory doesn't exist, create it
        if err := os.MkdirAll("dicom/out", os.ModePerm); err != nil {
            return "", err
        }
    }
    destinationDir = "./dicom/out"
    return destinationDir, nil
}

func checkContentType(expectedCt string, req *http.Request) error {
    ct := req.Header.Get("Content-Type")
    if ct == "" {
        return errors.New("Content-Type missing")
    }

    mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
    if mediaType != expectedCt {
        msg := "Content-Type header is not " + expectedCt
        return errors.New(msg)
    }

    return nil
}

type FileUploadRequest struct {
	File string `json:"file"`
}

func fileUploadHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != "POST" {
        http.Error(w, "Expected POST method", http.StatusBadRequest)
        return
    }

    // Parse the request body
    err := checkContentType("application/json", req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
        return
    }

    // Validate the supplied dataset
	reqBody, err := ioutil.ReadAll(req.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var uploadReq FileUploadRequest
	err = json.Unmarshal(reqBody, &uploadReq)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = dicom.ParseFile(uploadReq.File, nil)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Get the contents of the file to be uploaded
    fileContents, err := os.ReadFile(uploadReq.File)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Create a file on the server
    fileID := uuid.New().String()
    filePath := filepath.Join(destinationDir, fileID)
    dst, err := os.Create(filePath)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Copy the uploaded file data to the destination file
    //datasetJSON, err := json.Marshal(dataset)
    if err := os.WriteFile(filePath, fileContents, 0444); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Serve up the response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"fileID": fileID})
}

func main() {
    http.HandleFunc("/upload", fileUploadHandler)
    /*http.HandleFunc("/fileID/attributes", attributesHandler)
    http.HandleFunc("/fileID/png", pngHandler)*/
    _, err := createDestinationDir()
    if err != nil {
        log.Fatal(err.Error())
        return
    }
    log.Fatal(http.ListenAndServe(":8080", nil))
}
