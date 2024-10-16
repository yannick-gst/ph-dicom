package main

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"

	"github.com/google/uuid"
    "github.com/suyashkumar/dicom"
    "github.com/suyashkumar/dicom/pkg/tag"
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
    if err := os.WriteFile(filePath, fileContents, 0444); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Serve up the response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"fileID": fileID})
}

func attributesHandler(w http.ResponseWriter, req *http.Request) {
    if req.Method != "GET" {
        http.Error(w, "Expected GET method", http.StatusBadRequest)
        return
    }

    err := checkContentType("application/json", req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
        return
    }

    // Ensure the requested file exists and is a valid dataset
    parts := strings.Split(req.URL.Path, "/")
    fileID := parts[2]
    if fileID == "" {
        http.Error(w, "File ID is missing", http.StatusBadRequest)
        return
    }

	filePath := filepath.Join(destinationDir, fileID)
    if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
        http.Error(w, "File ID does not exist", http.StatusBadRequest)
        return
    }

    dataset, err := dicom.ParseFile(filePath, nil)
    if err != nil {
        http.Error(w, "A valid DICOM dataset is expected", http.StatusBadRequest)
        return
    }

    // Verify the tag requested
    tagGroup, err := strconv.ParseUint(req.URL.Query().Get("tagGroup"), 16, 16)
    if err != nil {
        http.Error(w, "Expected a valid tag group query parameter", http.StatusBadRequest)
        return
    }
    tagElement, err := strconv.ParseUint(req.URL.Query().Get("tagElement"), 16, 16)
    if err != nil {
        http.Error(w, "Expected a valid tag element query parameter", http.StatusBadRequest)
        return
    }

    // Lookup the element
    tag := tag.Tag{uint16(tagGroup), uint16(tagElement)}
    element, err := dataset.FindElementByTag(tag)
    if err != nil {
        http.Error(w, "The element could not be found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(element)
}

func main() {
    http.HandleFunc("/upload", fileUploadHandler)
    http.HandleFunc("/file/fileID/attributes", attributesHandler)
    //http.HandleFunc("/file/fileID/png", pngHandler)
    _, err := createDestinationDir()
    if err != nil {
        log.Fatal(err.Error())
        return
    }
    log.Fatal(http.ListenAndServe(":8080", nil))
}
