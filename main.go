package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
) 

const (
	ftpServer = "localhost:21"
	ftpUser = "one"
	ftpPass = "1234"
	remoteDir = "/ftp/one"
	localDir = "./"
)
func main() {

	conn, err := ftp.Dial(ftpServer, ftp.DialWithTimeout(5 * time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = conn.Login(ftpUser, ftpPass)
	if err != nil {
		log.Fatal(err)
	}

	syncDir(conn, localDir, remoteDir)

	fmt.Println("sent all the files to the server")
}

func syncDir(conn *ftp.ServerConn, localPath, remotePath string) {

	files, err := os.ReadDir(localPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		localFilePath := filepath.Join(localPath, file.Name())
		remoteFilePath := filepath.Join(remotePath, file.Name())

		if file.IsDir() {
			log.Println("sending a dir")

			if !dirExistsInRemote(conn, remoteFilePath) {
				err := conn.MakeDir(remoteFilePath)
				if err != nil {
					log.Fatalf("failed to create dir: %s %s", remoteFilePath, err)
				}
			}
		

			syncDir(conn, localFilePath, remoteFilePath)

		} else {
			localFile, err := os.Open(localFilePath)
			if err != nil {
				log.Fatal(err)
			}

			err = conn.Stor(remoteFilePath, localFile)
			if err != nil {
				log.Fatal(err)
			}

			localFile.Close()
		}

	}
}

func dirExistsInRemote(conn *ftp.ServerConn, remotePath string) bool {

	if err := conn.ChangeDir(remotePath); err != nil {
		return false
	}

	if err := conn.ChangeDir(".."); err != nil {
		return false
	}

	return true

}