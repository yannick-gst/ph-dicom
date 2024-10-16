package main

import (
	"bytes"
    "encoding/json"
	"fmt"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
	"net/http/httputil"
	"os"
	"path/filepath"
    "testing"

    //"github.com/suyashkumar/dicom"
)

type FileUploadResponse struct {
	FileID string `json:"fileID"`
}

func TestUploadHandler(t *testing.T) {
	createDestinationDir()

	var uploadReq FileUploadRequest
	cwd, _ := os.Getwd()
	uploadReq.File = filepath.Join(cwd, "testData", "IM000001")
	reqBody, err := json.Marshal(&uploadReq)
	if err != nil {
        t.Errorf("expected valid file, got %v", err)
	}

    req := httptest.NewRequest("POST", "/upload", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    fileUploadHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %v", w.Code)
	}

	res := w.Result()
	defer res.Body.Close()

    resBody, err := ioutil.ReadAll(res.Body)
    if err != nil {
        t.Errorf("expected error to be nil, got %v", err)
    }

	var uploadResp FileUploadResponse
	err = json.Unmarshal(resBody, &uploadResp)
	if err != nil {
        t.Errorf("expected response body in JSON format, got %v", err)
	}

	expectedFilePath := filepath.Join(destinationDir, uploadResp.FileID)
    _, err = os.Stat(expectedFilePath)
    if os.IsNotExist(err) {
        t.Errorf("expected a file to be uploaded, got %v", err)
    }
}

func TestAttributesHandler(t *testing.T) {
	createDestinationDir()

	// Ensure we have a file to test with
	cwd, _ := os.Getwd()
	srcFile := filepath.Join(cwd, "/testData", "18809e7d-ec3d-464e-b1bf-8073c50e8541")
	dstFile := filepath.Join(destinationDir, "18809e7d-ec3d-464e-b1bf-8073c50e8541")
	bytesRead, _ := ioutil.ReadFile(srcFile)
	ioutil.WriteFile(dstFile, bytesRead, 0644)

	url := "/file/18809e7d-ec3d-464e-b1bf-8073c50e8541/attributes?tagGroup=0010&tagElement=0010"
    req := httptest.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    attributesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %v", w.Code)
	}

	res := w.Result()
	defer res.Body.Close()

	respDump, _ := httputil.DumpResponse(res, true)
	fmt.Printf("RESPONSE:\n%s", string(respDump))
}