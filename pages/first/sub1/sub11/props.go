package sub11
import "net/http"
type Props struct {
	Num  int
	Name string
}

func GetProps(w http.ResponseWriter, req *http.Request) (Props,bool) {

	w.WriteHeader(http.StatusInternalServerError)

	return Props{Num: 1, Name: "Greg"},true

}
