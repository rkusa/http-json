package json

import (
	j "encoding/json"
	"net/http"
	"reflect"
	"strings"
)

// Read allows to read a JSON request body into a struct while filtering
// the read JSON properties by the provided whitelist.
func Read(req *http.Request, dst interface{}, whitelist []string) error {
	tmp := map[string]j.RawMessage{}

	body := req.Body
	defer body.Close()

	dec := j.NewDecoder(body)
	if err := dec.Decode(&tmp); err != nil {
		return err
	}

	fields := extract(dst)

	for _, name := range whitelist {
		field, ok := fields[name]
		if !ok {
			continue
		}

		raw, ok := tmp[name]
		if !ok {
			continue
		}

		val := reflect.New(field.Type())
		if err := j.Unmarshal(raw, val.Interface()); err != nil {
			return err
		}
		field.Set(val.Elem())
	}

	return nil
}

func extract(dst interface{}) map[string]reflect.Value {

	t := reflect.TypeOf(dst)
	v := reflect.ValueOf(dst)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	fields := map[string]reflect.Value{}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tags := strings.Split(f.Tag.Get("json"), ",")
		name := ""
		if len(tags) > 0 {
			name = tags[0]
		}
		if name == "-" {
			continue
		}
		if name == "" {
			name = f.Name
		}

		fv := v.Field(i)

		if f.Anonymous { // embedded struct
			ft := f.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				fv = fv.Elem()
			}

			if !fv.IsValid() { // eg. is nil
				// init embedded struct
				fv = reflect.New(ft)
				v.Field(i).Set(fv)
				fv = fv.Elem()
			}

			for k, v := range extract(fv.Addr().Interface()) {
				fields[k] = v
			}
		}

		fields[name] = fv
	}

	return fields
}
