package main

import (
	//"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/knq/ini"
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
	fileCfg, err := ini.LoadFile("config.cfg")
	if err != nil {
		log.Fatal("Error with service configuration %s", err)
	}

	port := fileCfg.GetKey("service-1.port")

	if port == "" {
		log.Fatal("Error with port number configuration")
	}

	gfService := CreateGFService()
	router := CreateNewRouter(gfService)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
