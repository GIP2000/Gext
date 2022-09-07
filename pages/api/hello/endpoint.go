package hello

type Something struct{
  Name string
  Value int
}

func Handle() Something {
  return Something{Name: "hi", Value: 12}
}
