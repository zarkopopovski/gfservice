package main

import (
	"crypto/sha1"
	"encoding/json"

	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"io"
	//"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	//"git.cerebralab.com/george/logo"
)

type FileController struct {
	dbManager *DBManager
}

func (fileManager *FileController) uploadFile(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	file, header, err := r.FormFile("file")

	userID := r.FormValue("user_id")
	parameter1 := r.FormValue("parameter_1")

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	defer file.Close()

	fileName := header.Filename

	randomFloat := strconv.FormatFloat(rand.Float64(), 'E', -1, 64)

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(time.Now().String() + fileName + userID + randomFloat))
	sha1HashString := sha1Hash.Sum(nil)

	fileNameHASH := fmt.Sprintf("%x", sha1HashString)

	fileName = fileNameHASH + fileName

	out, err := os.Create("./uploads/" + fileName)

	if err != nil {
		fmt.Fprintf(w, "Unable to create a file for writting. Check your write access privilege")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		fmt.Fprintln(w, err)
	}

	_ = fileManager.dbManager.addNewFile(userID, parameter1, fileName)

	fmt.Fprintf(w, userID+" File uploaded successfully : ")
	fmt.Fprintf(w, header.Filename)
}

func (fileManager *FileController) uploadMultipleFiles(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := r.FormValue("user_id")
	parameter1 := r.FormValue("parameter_1")

	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm

	files := m.File["files"]

	if len(files) > 0 {
		for i := range files {

			fileName := files[i].Filename

			randomFloat := strconv.FormatFloat(rand.Float64(), 'E', -1, 64)

			sha1Hash := sha1.New()
			sha1Hash.Write([]byte(time.Now().String() + fileName + userID + randomFloat))
			sha1HashString := sha1Hash.Sum(nil)

			fileNameHASH := fmt.Sprintf("%x", sha1HashString)

			fileName = fileNameHASH + fileName

			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			dst, err := os.Create("./uploads/" + fileName)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			_ = fileManager.dbManager.addNewFile(userID, parameter1, fileName)
		}
		fmt.Fprintf(w, userID+" Files uploaded successfully")
	} else {
		fmt.Fprintf(w, userID+" No Files for upload")
	}
}

func (fileManager *FileController) findUserFiles(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	userID := r.FormValue("user_id")
	parameter := r.FormValue("parameter_1")

	userFiles := []*FileModel{}

	if parameter == "" {
		userFiles = fileManager.dbManager.findAllUserFilesBy(userID)
	} else {
		userFiles = fileManager.dbManager.findUserFilesByUserIDAndParameter(userID, parameter)
	}

	if len(userFiles) > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(userFiles); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1"}); err != nil {
		panic(err)
	}
}
