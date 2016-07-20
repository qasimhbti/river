package river

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
)

type jsonDecoder []byte

func (j *jsonDecoder) Init(src io.Reader) error {
	b, err := ioutil.ReadAll(src)
	*j = b
	return err
}

func (j *jsonDecoder) Copy() []byte {
	buf := make([]byte, len(*j))
	copy(buf, *j)
	return buf
}

func (j *jsonDecoder) Decode(v interface{}) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			if _, ok := err1.(error); ok {
				err = err1.(error)
			} else {
				err = fmt.Errorf("Unitenfied error %v", err1)
			}
		}
	}()

	if !reflect.ValueOf(v).IsValid() || reflect.TypeOf(v).Kind() != reflect.Ptr || reflect.ValueOf(v).IsNil() {
		return fmt.Errorf("cannot marshal to %v, must be pointer and not nil", reflect.TypeOf(v))
	}

	if err := json.Unmarshal(j.Copy(), &v); err == nil {
		return nil
	} else if _, ok := err.(*json.UnmarshalTypeError); !ok {
		// not type related error
		return err
	}

	var elem reflect.Value

	if reflect.ValueOf(v).Elem().Kind() == reflect.Slice {
		// if type is slice, attempt an element of the slice
		// and return a slice containing the element.
		item := reflect.New(reflect.ValueOf(v).Elem().Type().Elem())
		if err := json.Unmarshal(j.Copy(), item.Interface()); err != nil {
			return err
		}
		elem = reflect.Append(reflect.ValueOf(v).Elem(), item.Elem())
	} else if reflect.ValueOf(v).Elem().Kind() == reflect.Struct {
		// if type is a struct, attempt a slice of the struct
		// and return first element of the slice.
		slice := reflect.New(reflect.SliceOf(reflect.TypeOf(v).Elem()))
		if err := json.Unmarshal(j.Copy(), slice.Interface()); err != nil {
			return err
		}
		elem = slice.Elem().Index(0)
	} else {
		return fmt.Errorf("Cannot unmarshal to type %v", reflect.ValueOf(v).Elem().Type())
	}

	reflect.ValueOf(v).Elem().Set(elem)
	return nil
}
