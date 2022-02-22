package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	EnvTag    = "env"
	PrefixTag = "env_prefix"
	Delimiter = ","
)

func OverwriteFromEnv(in interface{}) error {
	v := reflect.ValueOf(in)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("input must be pointer")
	}
	v = v.Elem()

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be pointer to struct")
	}
	t := v.Type()

	return traverse(v, t, "", "")
}

func traverse(v reflect.Value, t reflect.Type, tag reflect.StructTag, prefix string) (err error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(t.Elem()))
		}
		v = v.Elem()
		t = v.Type()
	}

	if v.Kind() == reflect.Struct {
		prefix = tag.Get(PrefixTag)
		for i := 0; i < v.NumField(); i++ {
			if err = traverse(v.Field(i), t.Field(i).Type, t.Field(i).Tag, prefix); err != nil {
				return err
			}
		}
		return
	}

	if !v.CanSet() {
		return
	}

	if k, ok := tag.Lookup(EnvTag); ok {
		k = prefix + k
		if val := os.Getenv(k); val != "" {
			if v.Kind() == reflect.Slice {
				if !isBasicType(t.Elem()) {
					return fmt.Errorf("not supported type: %v", t)
				}

				return setSlice(v, val)
			} else {
				if !isBasicType(t) {
					return fmt.Errorf("not supported type: %v", t)
				}

				return setBasicType(v, val)
			}
		}
	}

	return
}

var basicTypes = []reflect.Kind{
	reflect.Bool,
	reflect.Int,
	reflect.Int8,
	reflect.Int16,
	reflect.Int32,
	reflect.Int64,
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
	reflect.Float32,
	reflect.Float64,
	reflect.String,
}

func isBasicType(t reflect.Type) bool {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for _, tt := range basicTypes {
		if t.Kind() == tt {
			return true
		}
	}

	return false
}

func setSlice(v reflect.Value, val string) (err error) {
	raws := strings.Split(val, ",")
	ss := reflect.MakeSlice(v.Type(), len(raws), len(raws))
	for i := 0; i < len(raws); i++ {
		err = setBasicType(ss.Index(i), raws[i])
		if err != nil {
			return
		}
	}
	v.Set(ss)
	return
}

func setBasicType(v reflect.Value, val string) error {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return setBasicType(v.Elem(), val)
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
