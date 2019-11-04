package config

//go:generate go run ./generate.go

func GetExampleFile() string {
	return ExampleFile
}
