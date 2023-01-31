package main

//go:generate go run generate.go
import (
	"fmt"
	"log"
	"net/http"
	"Gext/routeMapper"
	"html/template"
)


type path struct {
	Path string
	Props string
}

func newPath (pth string, props string) path {
	if props == "" {
		props ="{}"
	}
	return path{Path:pth, Props:props}
}

func main() {

	t, err := template.ParseFiles("./public/index.html")

	if err != nil {
		log.Fatal("Please make sure you have a index.html file ")
		panic(err)
	}

	http.Handle("/static/",http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL.Path)
		if val, ok := routeMapper.RequestMap[req.URL.Path]; ok {
			if val.IsApi{
				v,earlyExit := val.HandleFunction(w,req)
				if earlyExit {
					return
				}
				w.Write(v)
				return ;
			} else {
				var props string = ""
				if val.HandleFunction != nil {
					initalProps, earlyExit := val.HandleFunction(w,req)
					if earlyExit {
						return
					}
					props = string(initalProps)
				}
				ptS := newPath(val.PathToTemplate, props)
				t.Execute(w,ptS)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
