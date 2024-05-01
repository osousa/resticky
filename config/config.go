package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/joho/godotenv"

	"github.com/osousa/resticky/pkg/sushi/toml"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func read[T any](ctx context.Context, c T) {
	// Load env file only if local env
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
		envpath := fmt.Sprintf("%s/../%s.env", basepath, env)
		if err := godotenv.Load(envpath); err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
	// Build config path to toml file
	f := fmt.Sprintf("%s/%s.toml", basepath, env)
	_, err := toml.DecodeFile(f, &c)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// Set and load custom config structs
func Read[T any](ctx context.Context, c T) T {
	// Check if c is a struct
	ptr := reflect.ValueOf(c)
	if ptr.Kind() != reflect.Ptr {
		panic("Value must be a pointer")
	}

	val := reflect.Indirect(ptr)
	if val.Kind() != reflect.Struct {
		panic("c must be a struct")
	}

	// Read config file
	read(ctx, c)

	// Iterate over fields of the struct
	for i := 0; i < val.NumField(); i++ {
		fieldInfo := val.Type().Field(i)
		field := val.Field(i)
		if field.IsZero() {
			panic(fmt.Sprintf("Field %s is not set", fieldInfo.Name))
		}
		fieldElem := field.Elem()
		if fieldElem.IsZero() {
			panic(fmt.Sprintf("Field %s missing config", fieldInfo.Name))
		}
	}
	return c
}
