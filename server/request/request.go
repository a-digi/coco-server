package request

import (
	"net/http"
	"reflect"
	"strconv"
	"sync"
)

var (
	formStructCache sync.Map
)

type formFieldInfo struct {
	Index   int
	JsonTag string
	Kind    reflect.Kind
}

func MapFormToStruct(r *http.Request, dest interface{}) error {
	v := reflect.ValueOf(dest)

	if v.Kind() != reflect.Ptr || v.IsNil() {
		return &MappingError{"dest must be a non-nil pointer to a struct"}
	}

	v = v.Elem()
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return &MappingError{"dest must be a pointer to a struct"}
	}

	var fields []formFieldInfo
	if cached, ok := formStructCache.Load(t); ok {
		fields = cached.([]formFieldInfo)
	} else {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			fields = append(fields, formFieldInfo{
				Index:   i,
				JsonTag: jsonTag,
				Kind:    field.Type.Kind(),
			})
		}
		formStructCache.Store(t, fields)
	}

	for _, info := range fields {
		formValue := r.FormValue(info.JsonTag)
		if formValue == "" {
			continue
		}

		fv := v.Field(info.Index)

		if !fv.CanSet() {
			continue
		}

		switch info.Kind {
		case reflect.String:
			fv.SetString(formValue)
		case reflect.Bool:
			b, err := strconv.ParseBool(formValue)
			if err == nil {
				fv.SetBool(b)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			iv, err := strconv.ParseInt(formValue, 10, 64)
			if err == nil {
				fv.SetInt(iv)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			uv, err := strconv.ParseUint(formValue, 10, 64)
			if err == nil {
				fv.SetUint(uv)
			}
		case reflect.Float32, reflect.Float64:
			fvf, err := strconv.ParseFloat(formValue, 64)
			if err == nil {
				fv.SetFloat(fvf)
			}
		}
	}
	return nil
}

type MappingError struct {
	msg string
}

func (e *MappingError) Error() string {
	return e.msg
}
