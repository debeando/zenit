package example

//go:generate go run ./generate.go

func GetConfigFile() string {
	return ExampleFile
}
