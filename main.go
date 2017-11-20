package main

import (
	//"flag"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type GFServCore struct {
	dbManager      *DBManager
	fileController *FileController
}

func CreateGFService() *GFServCore {
	gfServiceCore := &GFServCore{
		dbManager:      CreateDBConnection(),
		fileController: &FileController{},
	}

	gfServiceCore.fileController.dbManager = gfServiceCore.dbManager

	return gfServiceCore
}

func CreateNewRouter(handlers *GFServCore) *httprouter.Router {
	router := httprouter.New()

	router.POST("/upload_file", handlers.fileController.uploadFile)
	router.POST("/upload_files", handlers.fileController.uploadMultipleFiles)
	router.POST("/list_files", handlers.fileController.findUserFiles)

	router.ServeFiles("/uploads/*filepath", http.Dir("./uploads"))

	return router
}

func main() {
	gfService := CreateGFService()
	router := CreateNewRouter(gfService)

	log.Fatal(http.ListenAndServe(":8080", router))
}
