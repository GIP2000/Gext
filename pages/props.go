package pages
import "net/http"
type Props struct {
  Name string
}

func GetProps(w http.ResponseWriter, req *http.Request) (Props,bool){
  return Props{Name:"Greg"},false
}
