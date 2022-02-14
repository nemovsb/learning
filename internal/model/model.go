package model

type Task1Request struct {
	Key   string `json:"key"`
	Value int    `json:"val"`
}
type Task1Response struct {
	Key   string `json:"key"`
	Value int    `json:"val"`
}

type Task2RequestBody struct {
	StringValueForCheck string `json:"s"`
	SecretKey           string `json:"key"`
}

type Task3RequestBody [2]struct {
	A   string `json:"a"`
	B   string `json:"b"`
	Key string `json:"key"`
}
