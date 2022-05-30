package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/drive/v2"
	"log"
	"strings"
)

var totalFiles = 0

var totalDeletedFiles = 0

func main() {

	showInfo := false
	deleteFiles := true
	tokenFile := "Credentials/accessToken.json"
	credFile := "Credentials/desktopAppClient.json"
	folderId := "1JZ2zuyK3yJQrhsOHsTQUuusBvcMJj6SJ"
	driveScope := drive.DriveScope
	server := initAccessPoint(tokenFile, credFile, driveScope)
	loopOverFolder("", folderId, server, showInfo, deleteFiles)

	println()
	fmt.Printf("TOTAL FILES: %v\n", totalFiles)
	fmt.Printf("TOTAL DELETED FILES: %v\n", totalDeletedFiles)
}
func loopOverFolder(appendBefore string, folderID string, server *drive.Service, showInfo bool, deleteFiles bool) {

	elements, err := getAllChildren(server, folderID)

	if err != nil {
		log.Fatalf("cannot get all folders %v", err)
	}
	//_ = the element number, element= the actual object
	for num, element := range elements {
		fileInfo := getChildInfo(appendBefore, num, element, server, showInfo)
		isJson := checkIfJson(fileInfo)
		isFolder := checkIfFolder(fileInfo)

		if isJson {
			println(appendBefore + "JSON found!")
			println(appendBefore + "DELETED: " + fileInfo.Title)
			println()
			moveToTrash(server, fileInfo.Id, deleteFiles)
			totalDeletedFiles += 1
		} else if isFolder {
			println(appendBefore + "folder found!")

			loopOverFolder(appendBefore+"\t", element.Id, server, showInfo, deleteFiles)

		} else {
			println(appendBefore + "file found!")
		}

		totalFiles += 1

		if totalFiles%50 == 0 {
			fmt.Printf("TOTAL FILES: %v\n", totalFiles)
			fmt.Printf("TOTAL DELETED FILES: %v\n", totalDeletedFiles)
		}
	}

}

func moveToTrash(server *drive.Service, fileId string, deleteFiles bool) {

	if deleteFiles {

		_, err := server.Files.Trash(fileId).Do()
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
		}

	}

}

func checkIfFolder(info drive.File) bool {
	folderExtension := "/drive/folders/"
	title := info.AlternateLink
	if strings.Contains(title, folderExtension) {
		return true
	}
	return false
}

func checkIfJson(info drive.File) bool {
	jsonExtension := ".json"
	title := info.Title
	if strings.Contains(title, jsonExtension) {
		return true
	}
	return false
}

func getChildInfo(appendBefore string, num int, element *drive.ChildReference, server *drive.Service, showInfo bool) drive.File {

	info, _ := server.Files.Get(element.Id).Do()
	println(appendBefore + info.Title)
	if showInfo {
		fmt.Printf(appendBefore+"loc: %v\n", num)
		println(appendBefore + "id: " + element.Id)

		jsonString, _ := json.MarshalIndent(info, "", "    ")
		println(appendBefore + string(jsonString))
		println("===================================")

	}

	return *info
}
func getAllChildren(server *drive.Service, folderID string) ([]*drive.ChildReference, error) {
	var childRef []*drive.ChildReference
	pageToken := ""

	for {
		child := server.Children.List(folderID)

		if pageToken != "" {
			child = child.PageToken(pageToken)
		}

		resp, err := child.Do()

		if err != nil {
			fmt.Printf("an error occured: %v\n", err)
			return childRef, err
		}
		childRef = append(childRef, resp.Items...)

		pageToken = resp.NextPageToken
		if pageToken == "" {
			break
		}

	}
	return childRef, nil
}
