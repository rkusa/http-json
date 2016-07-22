// Package json provides helper methods to read (both filtered and unfiltered)
// and write JSON request bodies.
//
//  json.Read(r, &user)
//  json.ReadFiltered(r, &user, []string{"name"})
//  json.Write(rw, &user)
//
package json

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Write marshals the given data and writes it to the provided context.
func Write(rw http.ResponseWriter, data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.Write(json)
}

// Read reads the JSON body of the given context and decodes it into the
// provided data structure.
func Read(r *http.Request, v interface{}) error {
	body := r.Body
	defer body.Close()

	decoder := json.NewDecoder(body)
	err := decoder.Decode(v)
	switch err {
	case io.EOF:
		return nil
	default:
		return err
	}
}

// Read allows to read a JSON request body into a struct while filtering
// the read JSON properties by the provided whitelist.
func ReadFiltered(req *http.Request, dst interface{}, whitelist []string) error {
	tmp := map[string]json.RawMessage{}

	body := req.Body
	defer body.Close()

	dec := json.NewDecoder(body)
	err := dec.Decode(&tmp)
	if err != nil {
		switch err {
		case io.EOF:
			return nil
		default:
			return err
		}
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
		if err := json.Unmarshal(raw, val.Interface()); err != nil {
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
