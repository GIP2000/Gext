package sub2
import "net/http"
type Props struct {
	Num  int
	Name string
}

func GetProps(w http.ResponseWriter, req *http.Request) (Props, bool) {

	return Props{Num: 1, Name: "Greg"},false

}
