package cafebabe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
)

const fieldTagName = "cb"

func Unmarshal(data []byte, v any) error {
	r := bytes.NewReader(data)
	return NewDecoder(r).Decode(v)
}

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

func (d *Decoder) verifyDestination(v any) error {
	t := reflect.TypeOf(v)

	if t.Kind() != reflect.Ptr {
		return errors.New("dst must be a pointer")
	} else {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("*dst must be a struct, got %s", t)
	}

	return nil
}

func (d *Decoder) Decode(v any) error {
	if p, ok := v.(Primitive); ok {
		_, err := p.ReadFrom(d.r)
		return err
	}

	err := d.verifyDestination(v)
	if err != nil {
		return err
	}

	t := reflect.ValueOf(v).Elem()

	for i, vf := range reflect.VisibleFields(t.Type()) {
		if !vf.IsExported() || vf.Tag.Get(fieldTagName) == "-" {
			continue
		}

		field := t.FieldByIndex(vf.Index)

		if vf.Tag.Get(fieldTagName) == "variadic" {
			prepareMethodName := fmt.Sprintf("Prepare%s", vf.Name)
			m := t.Addr().MethodByName(prepareMethodName)

			if !m.IsValid() {
				return fmt.Errorf("%s should have %s method in order to process variadic field %s", t, m, vf.Name)
			}

			res := m.Call([]reflect.Value{})
			if len(res) != 1 {
				return fmt.Errorf("%s should have exactly one returned value, got %d", prepareMethodName, len(res))
			}

			if !res[0].CanInterface() {
				return fmt.Errorf("invalid value returned from %s", prepareMethodName)
			}

			if it, ok := res[0].Interface().(reflect.Type); ok {
				v := reflect.New(it)
				err := d.Decode(v.Interface())
				if err != nil {
					return fmt.Errorf("error decoding variadic field %s: %s", vf.Name, err)
				}

				field.Set(v.Elem())
				continue
			} else {
				return fmt.Errorf("invalid value returned from %s: expected reflect.Type, got %s", prepareMethodName, it)
			}
		}

		v := field.Addr().Interface()
		if p, ok := v.(Primitive); ok {
			_, err := p.ReadFrom(d.r)
			if err != nil {
				return err
			}

			continue
		}

		if field.Kind() == reflect.Slice {
			if i == 0 {
				return errors.New("array could not be the first item in struct")
			}

			prev := t.Field(i - 1)
			size := int(prev.Uint())

			if vf.Tag.Get(fieldTagName) == "1-indexed" {
				size -= 2 // TODO
			}

			it := field.Type().Elem()
			s := reflect.MakeSlice(reflect.SliceOf(it), size, size)

			for j := 0; j < size; j++ {
				iv := reflect.New(it)

				err := d.Decode(iv.Interface())
				if err != nil {
					return fmt.Errorf("field %s: error decoding element %d: %s", vf.Name, j, err)
				}

				s.Index(j).Set(iv.Elem())
			}

			field.Set(s)
			continue
		}

		if field.Kind() == reflect.Struct {
			err := d.Decode(field.Addr().Interface())
			if err != nil {
				return fmt.Errorf("field %s: error decoding struct: %s", vf.Name, err)
			}

			continue
		}

		return fmt.Errorf("cannot decode field %s (%s)", vf.Name, field.Type())
	}

	return nil
}
