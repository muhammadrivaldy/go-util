package goutil

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"
)

// Configuration is a function for get info configuration
func Configuration(osFile *os.File, model interface{}) error {

	decoder := json.NewDecoder(osFile)
	err := decoder.Decode(model)
	if err != nil {
		return err
	}

	if err := envConfiguration(model); err != nil {
		return err
	}

	return err

}

func envConfiguration(req interface{}) error {

	v := reflect.ValueOf(req)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return setEnv(v.Type(), v)

}

func setEnv(t reflect.Type, v reflect.Value) error {

	if v.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := v.Field(i)
			k := v.Field(i).Kind()

			if k == reflect.Struct {
				if err := setEnv(f.Type(), f); err != nil {
					return err
				}

				continue
			}

			val := os.Getenv(t.Field(i).Tag.Get("env"))
			if val == "" {
				continue
			}

			switch k {
			case reflect.String:
				f.SetString(val)
			case reflect.Int:
				valInt, err := strconv.ParseInt(val, 10, 0)
				if err != nil {
					return err
				}
				f.SetInt(valInt)
			case reflect.Int8:
				valInt, err := strconv.ParseInt(val, 10, 8)
				if err != nil {
					return err
				}
				f.SetInt(valInt)
			case reflect.Int16:
				valInt, err := strconv.ParseInt(val, 10, 16)
				if err != nil {
					return err
				}
				f.SetInt(valInt)
			case reflect.Int32:
				valInt, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					return err
				}
				f.SetInt(valInt)
			case reflect.Int64:
				valInt, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}
				f.SetInt(valInt)
			case reflect.Bool:
				valBool, err := strconv.ParseBool(val)
				if err != nil {
					return err
				}
				f.SetBool(valBool)
			case reflect.Float32:
				valFloat, err := strconv.ParseFloat(val, 32)
				if err != nil {
					return err
				}
				f.SetFloat(valFloat)
			case reflect.Float64:
				valFloat, err := strconv.ParseFloat(val, 64)
				if err != nil {
					return err
				}
				f.SetFloat(valFloat)
			}
		}
	}

	return nil

}
