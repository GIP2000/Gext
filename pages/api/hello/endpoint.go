package hello

import "net/http"

type Something struct{
  Name string
  Value int
}

func Handle(w http.ResponseWriter, req *http.Request) (Something,bool) {

  return Something{Name: "hi", Value: 12},false
}
