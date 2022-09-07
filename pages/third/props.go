package third

type Props struct {
	Num  int
	Name string
}

func GetProps() Props {

	return Props{Num: 1, Name: "Greg"}

}
