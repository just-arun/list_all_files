package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
)

// Item struct
type Item struct {
	Name string `json:"name"`
}

var items []Item
var dirName string
var outputDir string

func createJSON() []byte {
	data, err := json.Marshal(items)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func populateList() {
	if len(dirName) < 1 {
		dirName = "."
	}

	fmt.Println("pop", dirName)

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		name := f.Name()
		items = append(items, Item{Name: name})
	}
}

func setHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	return w
}

func getItems(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")

	if len(path) > 0 {
		dirName = path
	}

	fmt.Println(dirName)
	populateList()
	data := createJSON()
	wr := setHeader(w)
	wr.Write(data)
}

func main() {

	fmt.Printf("Enter the directive (by default it is current directive): ")
	fmt.Scanln(&dirName)

	populateList()

	fmt.Printf("Options \n 1) json file \n 2) server\nEnter Your output type:")

	var outputType int
	fmt.Scanln(&outputType)
	if outputType == 1 || outputType != 2 {
		fmt.Printf("Enter the output Directive: ")

		fmt.Scanln(&outputDir)

		if len(outputDir) < 1 {
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			outputDir = dir
		}
		date1 := time.Now()

		data := createJSON()
		err := ioutil.WriteFile(outputDir+"/data.json", data, 0644)
		if err != nil {
			log.Fatal(err)
		}
		date2 := time.Now()
		dureation := date2.Sub(date1).Seconds()
		fmt.Println("file generated in ", dureation, "seconds")
	} else {
		var port string
		fmt.Print("Enter the port to run the server:")
		fmt.Scanln(&port)
		matched, _ := regexp.MatchString(`\d{4}`, port)
		fmt.Println(matched)
		if !matched {
			log.Fatal("You have entered a invalid port number")
			port = "3300"
		}
		http.HandleFunc("/items", getItems)
		log.Fatal(http.ListenAndServe(":"+port, nil), "server started at http://localhost:3300")
	}

}
