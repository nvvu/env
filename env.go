package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const EnvTag = "env"

func OverwriteFromEnv(in interface{}) error {
	v := reflect.ValueOf(in)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("input must be pointer")
	}

	v = v.Elem()

	return inspect(v, "")
}

func inspect(v reflect.Value, tag reflect.StructTag) (err error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			if err = inspect(v.Field(i), t.Field(i).Tag); err != nil {
				return err
			}
		}
		return
	}

	if !v.CanSet() {
		return
	}

	if k, ok := tag.Lookup(EnvTag); ok {
		if val := os.Getenv(k); val != "" {
			return setPrimitiveType(v, val)
		}
	}

	return
}

func setPrimitiveType(v reflect.Value, val string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(val)
	case reflect.Bool:
		i, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		v.SetBool(i)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 0, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(val, 0, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		v.SetFloat(i)
	}
	return nil
}

func typeInfo(t reflect.Type) {
	fmt.Println("type", t)
	fmt.Println("> kind:", t.Kind().String())
	fmt.Println("> comparable:", t.Comparable())
	fmt.Println("> number of methods:", t.NumMethod())

	if t.Kind() == reflect.Struct {
		fmt.Println("> number of field:", t.NumField())
		for i := 0; i < t.NumField(); i++ {
			fmt.Println(">> field: ", i, t.Field(i), t.Field(i).Name)
		}
	}
	fmt.Println("")
}

func valueInfo(t reflect.Value) {
	fmt.Println("value", t, t.Type())
	fmt.Println("> kind:", t.Kind().String())
	fmt.Println("> can addresable:", t.CanAddr())
	fmt.Println("> can set:", t.CanSet())
	fmt.Println("> number of methods:", t.NumMethod())
	fmt.Println(t.Type())

	if t.Kind() == reflect.Struct {
		fmt.Println("> number of field:", t.NumField())
		for i := 0; i < t.NumField(); i++ {
			vf := t.Field(i)

			fmt.Println(">> field: ", i, vf, vf.Type())
			fmt.Println(vf.CanSet(), vf.CanAddr())
			if vf.Kind() == reflect.Struct {
				valueInfo(vf)
			}

		}
	}

	fmt.Println("")
}
