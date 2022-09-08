package badRequest

import "net/http"

type BR struct {
  Status int32
}


func Handle(w http.ResponseWriter, req *http.Request) (BR,bool) {
  w.WriteHeader(http.StatusInternalServerError)
  return BR{Status:200},true
}
