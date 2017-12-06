package encode

import (
	json "encoding/json"
	"reflect"
)

type DefaultSerializer struct {
	generic reflect.Value
}

/*
 * Creates a "generic" Serializer. v must be a 0 value of the desired type.
 */
func NewDefaultSerializer(v interface{}) *DefaultSerializer {
	s := new(DefaultSerializer)
	s.generic = reflect.ValueOf(v)

	return s
}

func (this *DefaultSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (this *DefaultSerializer) Deserialize(v []byte) (interface{}, error) {
	data := reflect.Indirect(this.generic).Interface()
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}

	return data, nil
}
