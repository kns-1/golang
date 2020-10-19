// A dynamic Golang program that parses files in the Bitbucket repository and validates if the file text structure matches the requirements of the repository structure or not.

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type YamlStruct struct {
	SkillsMap map[string]interface{} `yaml:"skills,omitempty"`
}

func dirExists(dirKey string, rootpath string, dirValue interface{}) bool {

	//	fmt.Println("KEY = ", dirKey)
	tempPath := rootpath + "/" + dirKey

	if nil == dirValue {
		//	fmt.Println("DETECTED")
		_, err := ioutil.ReadDir(tempPath)
		if err != nil {
			//log.Fatal(err)
			fmt.Println("PARENT KEY = ", dirKey)
			return false
		}
		return true
	}

	//fmt.Println(tempPath)
	files, err := ioutil.ReadDir(tempPath)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("PARENT KEY = ", dirKey)
		return false
	}

	flag := 0

	for _, listValues := range dirValue.([]interface{}) {
		//fmt.Println(listValues)
		for key := range listValues.(map[interface{}]interface{}) {
			//fmt.Println(key)
			flag = 0
			for _, fileNameObj := range files {
				//fmt.Println(f.Name())
				//removing filename extension for comparision
				fileName := fileNameObj.Name()
				filePrefix := strings.TrimSuffix(fileName, filepath.Ext(fileName))
				//fmt.Println(filePrefix)

				//fmt.Println("PREFIX = ", strings.ToLower(filePrefix))

				//convert interface{} to string type
				myKey := fmt.Sprintf("%v", key)

				//fmt.Println("MYKEY = ", strings.ToLower(myKey))

				eq := reflect.DeepEqual(strings.ToLower(filePrefix), strings.ToLower(myKey))
				if eq {
					//fmt.Println("EQUAL")
					flag = 1
					break
				}
			}
			if 1 == flag {
				break
			}
		}
		if 0 == flag {
			fmt.Println("PARENT KEY = ", dirKey)
			fmt.Println(listValues)
			return false
		}
	}
	return true
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	var filePath string
	var retDirValue bool
	fmt.Println("Enter root path: ")
	_ = scanner.Scan()
	filePath = scanner.Text()
	rootpath := filePath + "/skills"
	//fmt.Println("Root: ", rootpath)

	pathYaml := filePath + "/myFolder"

	yamlFiles, err := ioutil.ReadDir(pathYaml)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileNameObj := range yamlFiles {
		//fmt.Println(fileNameObj.Name())

		if ".yml" == filepath.Ext(fileNameObj.Name()) {

			pathYamlFile := pathYaml + "/" + fileNameObj.Name()
			fmt.Println("Path: ", pathYamlFile)

			parseYaml, parseYamlErr := ioutil.ReadFile(pathYamlFile)
			if parseYamlErr != nil {
				fmt.Errorf("error reading yaml file: %v", parseYamlErr)
			}

			var yamlStructObj YamlStruct
			err = yaml.Unmarshal([]byte(parseYaml), &yamlStructObj)
			if err != nil {
				fmt.Printf("error %+v", err)
			}
			for dirKey, dirValue := range yamlStructObj.SkillsMap {
				//	fmt.Printf("type=%+v key=%+v type=%+v value=%+v\n", reflect.TypeOf(dirKey), dirKey, reflect.TypeOf(dirValue), dirValue)

				retDirValue = dirExists(dirKey, rootpath, dirValue)
				if false == retDirValue {
					log.Println("KEY DOES NOT EXIST")
				}

			}
		}
	}
}
