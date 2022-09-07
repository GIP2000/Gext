//go:build exclude

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type TemplateEndpoint struct {
	Path         string
	TemplatePath string
	PropsFunc    bool
}

var templates []TemplateEndpoint
var apiEndpoints map[string]int
var apiEnpointCount int = 0

func removePages(path string) string {
	index := strings.Index(path, "pages")
	if index < 0 {
		panic("Invalid Path")
	}
	newStr := path[index+5:]
	if newStr == "" {
		return "/"
	}
	return newStr
}

func walkFunction(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		return nil
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	template := TemplateEndpoint{Path: path, PropsFunc: false}

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(path,"pages/src"){
			continue
		}

		if strings.HasPrefix(path, "pages/api") && strings.Contains(file.Name(), ".go") {
			if apiEndpoints == nil {
				apiEndpoints = make(map[string]int)
			}
			apiEndpoints[path] = apiEnpointCount
			apiEnpointCount += 1
			continue
		}

		if strings.Contains(file.Name(), ".svelte") {
			template.TemplatePath = path + "/" + file.Name()
		}
		template.PropsFunc = template.PropsFunc || strings.Contains(file.Name(), ".go")

	}

	if template.TemplatePath != "" {
		templates = append(templates, template)
	}

	return nil
}


func makeFile (){

	parent := "Gext" // TODO: Make this get the actualy package name

	f, err := os.OpenFile("routeMapper/routesMap.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 411)
	if err != nil {
		panic(err)
	}


	imports := ""
	for path, i := range apiEndpoints {
		imports += "apiEndpoint" + strconv.Itoa(i) + " \"" + parent + "/" + path + "\"\n"
	}
	for i, template := range templates {
		if !template.PropsFunc {
			continue
		}
		imports += "routeEndpoint" + strconv.Itoa(i) + " \"" + parent + "/" + template.Path + "\"\n"
	}
	mapStr := "map[string]Endpoint{"

	for i, template := range templates {
		mapStr += "\"" + removePages(template.Path) + "\": {PathToTemplate:\"" + "./pubic/routeEndpoint" + strconv.Itoa(i) + "/bundle.js" + "\", IsApi: false"
		if template.PropsFunc {
			mapStr += ",HandleFunction: func() []byte {s,err := json.Marshal(routeEndpoint" + strconv.Itoa(i) + ".GetProps()); if err != nil {panic (err)}; return s;}"
		}
		mapStr += "}"
		if i < len(templates)-1 || len(apiEndpoints) > 0 {
			mapStr += ","
		}
	}

	for path, i := range apiEndpoints {
		mapStr += "\"" + removePages(path) + "\": {IsApi: true, HandleFunction: func () []byte {s,err := json.Marshal(apiEndpoint" + strconv.Itoa(i) + ".Handle()); if err != nil {panic(err)}; return s;}}"
		if i < len(apiEndpoints)-1 {
			mapStr += ","
		}
	}

	mapStr += "}"

	_, err = f.WriteString(fmt.Sprintf("package routeMapper\n import(\n%s\n\"encoding/json\"\n)\ntype Endpoint struct {\nPathToTemplate string\nHandleFunction func() []byte\nIsApi bool\n}\n var RequestMap map[string]Endpoint = %s", imports, mapStr))
	if err != nil {
		panic(err)
	}
}
func bundleAll() {
	names, err := ioutil.ReadDir("./public")
	if err != nil {
		panic (err)
	}

	for _,entry := range names {

		if entry.IsDir() {
			os.RemoveAll("./public/" + entry.Name())
		}
	}

	for i,template :=range templates {
			bundler(template.TemplatePath,i)
		}
}

func bundler(path string, index int) {
	f,err := os.OpenFile("src/main.ts", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 411)	
	if err != nil {
		panic(err)
	}

	f.WriteString("//This was Generated\nimport App from \"." + path + "\"")
	f.Close()

	cmd := exec.Command("npm","run","build")
	dir,err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fmt.Println("dir", dir)
	
	cmd.Dir = dir
	
	out,err := cmd.Output()

	fmt.Println(string(out))
	
	if err != nil {
		panic(err)
	}
	
	err = os.Rename(dir + "/public/build/", dir + "/public/routeEndpoint" + strconv.Itoa(index))
	if err != nil {
		panic(err)
	}	

}

func main() {
	root := "./pages" // TODO: Make sure this is based off the file the generate comment is in 
	err := filepath.Walk(root, walkFunction)
	if err != nil {
		panic(err)
	}
	bundleAll()
	// TODO: Build Svelte Applications and save the bundle paths
	makeFile()


}
