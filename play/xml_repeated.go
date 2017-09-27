package play

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// func (s *SomeSlice) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
// }

func marshalRepeated(value interface{}, e *xml.Encoder, start xml.StartElement) error {
	v := reflect.ValueOf(value)
	t := v.Type()
unpack:
	for {
		switch t.Kind() {
		case reflect.Ptr, reflect.Interface:
			v = v.Elem()
			t = t.Elem()
		default:
			break unpack
		}
	}
	if t.Kind() != reflect.Slice {
		return errors.New("not a slice")
	}
	tf := t.Elem()

	m := new(xml.Marshaler)
	marshaler := tf.Implements(reflect.ValueOf(m).Elem().Type())

	for i := 0; i < v.Len(); i++ {
		vf := v.Index(i)
		s := xml.StartElement{
			Attr: start.Attr,
			Name: xml.Name{
				Space: start.Name.Space,
				Local: fmt.Sprintf("%s_%d", start.Name.Local, i),
			},
		}
		if marshaler {
			ret := vf.MethodByName("MarshalXML").Call([]reflect.Value{reflect.ValueOf(e), reflect.ValueOf(s)})
			if len(ret) > 0 {
				if err, ok := ret[0].Interface().(error); ok && err != nil {
					return err
				}
			}
		} else {
			if err := e.EncodeElement(vf.Interface(), s); err != nil {
				return err
			}
		}
	}

	return nil
}

func unmarshalRepeated(value interface{}, d *xml.Decoder, start xml.StartElement) error {
	v := reflect.ValueOf(value)
	t := v.Type()
unpack:
	for {
		switch t.Kind() {
		case reflect.Ptr, reflect.Interface:
			v = v.Elem()
			t = t.Elem()
		default:
			break unpack
		}
	}
	if t.Kind() != reflect.Slice {
		return errors.New("not a slice")
	}
	tf := t.Elem()

	vf := reflect.New(tf)

	name := strings.TrimRightFunc(start.Name.Local, unicode.IsDigit)
	name = strings.TrimSuffix(name, "_")

	s := xml.StartElement{
		Attr: start.Attr,
		Name: xml.Name{
			Space: start.Name.Space,
			Local: name,
		},
	}

	u := new(xml.Unmarshaler)
	if tf.Implements(reflect.ValueOf(u).Elem().Type()) {
		ret := vf.MethodByName("UnmarshalXML").Call([]reflect.Value{reflect.ValueOf(d), reflect.ValueOf(&s)})
		if len(ret) > 0 {
			if err, ok := ret[0].Interface().(error); ok && err != nil {
				return err
			}
		}
	} else {
		if err := d.DecodeElement(vf.Interface(), &s); err != nil {
			return err
		}
	}

	if v.Cap() < v.Len()+1 {
		if v.Cap() == 0 {
			v.SetCap(1)
		} else {
			v.SetCap(2 * v.Cap())
		}
	}
	v.SetLen(v.Len() + 1)
	v.Index(v.Len() - 1).Set(vf)
	return nil
}

func unmarshalElement(value interface{}, d *xml.Decoder, start xml.StartElement) error {
	val := reflect.ValueOf(value)
	typ := val.Type()

	name := strings.TrimRightFunc(start.Name.Local, unicode.IsDigit)
	name = strings.TrimSuffix(name, "_")
	s := xml.StartElement{
		Attr: start.Attr,
		Name: xml.Name{
			Space: start.Name.Space,
			Local: name,
		},
	}

	u := new(xml.Unmarshaler)
	if typ.Implements(reflect.ValueOf(u).Elem().Type()) {
		ret := val.MethodByName("UnmarshalXML").Call([]reflect.Value{reflect.ValueOf(d), reflect.ValueOf(&s)})
		if len(ret) > 0 {
			if err, ok := ret[0].Interface().(error); ok && err != nil {
				return err
			}
		}
	} else {
		if err := d.DecodeElement(val.Interface(), &s); err != nil {
			return err
		}
	}
	return nil
}
