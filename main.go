package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type JsMap struct {
	Version        int      `json:"version"`
	Sources        []string `json:"sources"`
	Names          []string `json:"names"`
	Mappings       string   `json:"mappings"`
	File           string   `json:"file"`
	SourcesContent []string `json:"sourcesContent"`
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Map file not found")
		fmt.Println("Specify by Args")
		return
	}

	arg := args[1]
	if !strings.Contains(arg, "\\") {
		arg = GetCurrentDir() + "\\" + arg
	}

	byte, err := ioutil.ReadFile(arg)
	if err != nil {
		panic(err)
	}

	var jsMap JsMap
	err = json.Unmarshal(byte, &jsMap)
	if err != nil {
		panic(err)
	}

	for i, source := range jsMap.SourcesContent {
		lines := strings.Split(source, "\n")
		output := fmt.Sprintf(
			"%s\\%s",
			GetCurrentDir(),
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						strings.ReplaceAll(
							strings.ReplaceAll(
								jsMap.Sources[i],
								"webpack://", "",
							),
							" ", "_",
						),
						"./", "",
					),
					"?", "",
				),
				"\u0000#/", "",
			),
		)

		if strings.Contains(output, "/") {
			ps := strings.Split(output, "/")
			slice := ps[0 : len(ps)-1]
			var dir string
			for _, str := range slice {
				dir += str + "/"
			}

			println("DIR: " + dir)
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}

		if !strings.Contains(output, ".js") && !strings.Contains(output, ".ts") {
			output += ".js"
		}

		_, err = os.Create(output)
		if err != nil {
			panic(err)
		}

		outputFile, err := os.OpenFile(output, os.O_RDWR|os.O_APPEND, 0660)
		if err != nil {
			panic(err)
		}

		for _, src := range lines {
			src = strings.ReplaceAll(src, "\t", "	")
			_, err = outputFile.WriteString(src + "\n")
			if err != nil {
				log.Fatal(err)
			}
		}

		err = outputFile.Close()
		if err != nil {
			panic(err)
		}
	}
}

func GetCurrentDir() string {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(path)
	return exPath
}
