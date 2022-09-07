package main

//go:generate go run generate.go
import (
	"fmt"
	"log"
	"net/http"
	"Gext/routeMapper"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL.Path)
		if val, ok := routeMapper.RequestMap[req.URL.Path]; ok {
			if val.IsApi{
				w.Write(val.HandleFunction())
			} else {
				w.Write([]byte(val.PathToTemplate))
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}

	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
