package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"crypto/sha1"
	"time"
)

type DBManager struct {
	db *sql.DB
}

type FileModel struct {
	Id           string `json:"id"`
	UserID       string `json:"user_id"`
	Parameter    string `json:"parameter"`
	FileName     string `json:"file_name"`
	Deleted      string `json:"deleted"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`
}

func CreateDBConnection() *DBManager {
	dbManager := &DBManager{}
	dbManager.OpenConnection()

	return dbManager
}

func (dbConnection *DBManager) OpenConnection() (err error) {
	db, err := sql.Open("sqlite3", "./gfserv.db")
	if err != nil {
		panic(err)
	}

	fmt.Println("SQLite Connection is Active")
	dbConnection.db = db

	dbConnection.setupInitialDatabase()

	return
}

func (dbConnection *DBManager) setupInitialDatabase() (err error) {
	statement, _ := dbConnection.db.Prepare("CREATE TABLE IF NOT EXISTS files (id VARCHAR PRIMARY KEY, user_id VARCHAR, parameter1 VARCHAR, file_name VARCHAR, deleted INTEGER, date_created VARCHAR, date_modified VARCHAR)")
	statement.Exec()

	return
}

func (dbConnection *DBManager) addNewFile(userID string, parameter1 string, fileName string) (err error) {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(time.Now().String() + userID + parameter1 + fileName))
	sha1HashString := sha1Hash.Sum(nil)

	fileID := fmt.Sprintf("%x", sha1HashString)

	query := "INSERT INTO files(id, user_id, parameter1, file_name, deleted, date_created, date_modified) VALUES($1, $2, $3, $4, 0, datetime('now'), datetime('now'))"

	_, err = dbConnection.db.Exec(query, fileID, userID, parameter1, fileName)

	if err == nil {
		return nil
	}

	panic(err)

	return
}

func (dbConnection *DBManager) markFileAsDeleted(userID string, fileID string) (err error) {
	query := "UPDATE files SET deleted=1, date_modified=datetime('now') WHERE id=$1 AND user_id=$2"

	_, err = dbConnection.db.Exec(query, fileID, userID)

	if err == nil {
		return nil
	}

	panic(err)

	return
}

func (dbConnection *DBManager) findAllUserFilesBy(userID string) []*FileModel {
	query := "SELECT id, user_id, parameter1, file_name, deleted, date_created FROM files WHERE user_id=$1 ORDER BY date_created DESC"

	rows, err := dbConnection.db.Query(query, userID)

	if err == nil {
		filesModels := make([]*FileModel, 0)

		for rows.Next() {
			newFile := new(FileModel)

			_ = rows.Scan(&newFile.Id, &newFile.UserID, &newFile.Parameter, &newFile.FileName, &newFile.Deleted, &newFile.DateCreated)

			filesModels = append(filesModels, newFile)
		}

		return filesModels
	}

	return nil
}

func (dbConnection *DBManager) findUserFilesByUserIDAndParameter(userID string, parameter string) []*FileModel {
	query := "SELECT id, user_id, parameter1, file_name, deleted, date_created FROM files WHERE user_id=$1 AND parameter1=$2 ORDER BY date_created DESC"

	rows, err := dbConnection.db.Query(query, userID, parameter)

	if err == nil {
		filesModels := make([]*FileModel, 0)

		for rows.Next() {
			newFile := new(FileModel)

			_ = rows.Scan(&newFile.Id, &newFile.UserID, &newFile.Parameter, &newFile.FileName, &newFile.Deleted, &newFile.DateCreated)

			filesModels = append(filesModels, newFile)
		}

		return filesModels
	}

	return nil
}
