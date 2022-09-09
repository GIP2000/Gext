//go:build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"regexp"
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

func makeFunc(name string) string {
	return "func(w http.ResponseWriter, req *http.Request) ([]byte,bool) {val, exitEarly := " + name + "(w,req); if exitEarly { return nil, true}; s,err := json.Marshal(val); if err != nil {panic(err)}; return s,false}"
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
		if strings.HasPrefix(template.Path,"./") {
			template.Path = template.Path[2:]
		}
		imports += "routeEndpoint" + strconv.Itoa(i) + " \"" + parent + "/" + template.Path + "\"\n"
	}
	mapStr := "map[string]Endpoint{"

	for i, template := range templates {
		mapStr += "\"" + removePages(template.Path) + "\": {PathToTemplate:\"" + "routeEndpoint" + strconv.Itoa(i) + "\", IsApi: false"
		if template.PropsFunc {
			mapStr += ",HandleFunction: " + makeFunc("routeEndpoint" + strconv.Itoa(i) + ".GetProps")
		}
		mapStr += "}"
		if i < len(templates)-1 || len(apiEndpoints) > 0 {
			mapStr += ","
		}
	}

	for path, i := range apiEndpoints {
		mapStr += "\"" + removePages(path) + "\": {IsApi: true, HandleFunction:" + makeFunc("apiEndpoint" + strconv.Itoa(i) + ".Handle") + "}"
		if i < len(apiEndpoints)-1 {
			mapStr += ","
		}
	}

	mapStr += "}"

	_, err = f.WriteString(fmt.Sprintf("package routeMapper\n import(\n%s\n\"encoding/json\"\n\"net/http\"\n)\ntype Endpoint struct {\nPathToTemplate string\nHandleFunction func(http.ResponseWriter,*http.Request) ([]byte,bool)\nIsApi bool\n}\n var RequestMap map[string]Endpoint = %s", imports, mapStr))
	if err != nil {
		panic(err)
	}
}

func makeTSFiles () {
	parent := "Gext"
	for path, _ := range apiEndpoints {
		makeTSFile(parent,path,path + "/endpoint.go",true)
	}

	for _,template :=range templates {
		path := template.Path
		fmt.Println("trying " + path)
		if strings.HasPrefix(path,".") || strings.HasPrefix(path,"/") {
			path = path[1:]
		}
		if strings.HasPrefix(path,"/") {
			path = path[1:]
		}
		makeTSFile(parent,path, path + "/props.go",false)
	}
}

func makeTSFile (parent string, packagePath string, scriptPath string, addFetch bool) {
	fileName := "types.ts"
	if addFetch {
		fileName = "request.ts"
	}
	cmd := exec.Command("./tscriptify","-interface", "-package=" + parent + "/" + packagePath, "-target="  + packagePath + "/" + fileName, scriptPath) 
	dir, err := os.Getwd()

	if err != nil {
		panic (err)
	}
	cmd.Dir = dir

	out,err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		panic(err)
	}
	
	if !addFetch {
		return 
	}

	fileWrite,err := os.OpenFile(dir + "/" + packagePath + "/" + fileName, os.O_APPEND | os.O_WRONLY, 411)

	if err != nil {
		panic (err)
	}

	defer fileWrite.Close()

	scriptBytes,err := os.ReadFile(dir + "/" + scriptPath)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(` *func *Handle *\(.*\) *\((.*) *, *bool *\)`)
	matches := re.FindSubmatch(scriptBytes)
	if len(matches) < 1 {
		fmt.Println("Couldn't find handle function")
		return ;
	}
	DataType := string(matches[1])
	fmt.Println(DataType)

	if _, err = fileWrite.WriteString("\nexport const getData = async ():Promise<" + DataType + "> =>(await fetch(\"" + removePages(packagePath) + "\")).json() as Promise<" + DataType + ">"); err != nil {
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

	if !strings.HasPrefix(path,"./") {
		path = "./" + path
	}

	f.WriteString("//This was Generated\nimport App from \"." + path + "\";\nexport default App;\n")
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
	makeTSFiles()
	bundleAll()
	makeFile()


}
