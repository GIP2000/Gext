package lib

import "reflect"

type apiReturnType[T any] struct {
  EarlyReturn bool
  Data T
}

type defaultData[T any] struct {
  Data T
}

func MakeReturnType[T any](data T, earlyReturn bool) apiReturnType[T]{
  rValue := reflect.ValueOf(data)
  if rValue.Kind() == reflect.Struct {
    return apiReturnType[T]{Data:data,EarlyReturn:earlyReturn}
  }
  return apiReturnType[T]{EarlyReturn: earlyReturn}
}

