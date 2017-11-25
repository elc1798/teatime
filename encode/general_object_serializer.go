package encode

import (
	json "encoding/json"
	"log"
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

	log.Printf("DefaultSerializer: Got type %v", s.generic)
	return s
}

func (this *DefaultSerializer) Serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (this *DefaultSerializer) Deserialize(v []byte) (interface{}, error) {
	data := reflect.Indirect(this.generic).Interface()
	log.Printf("DefaultSerializer: Using type %v", reflect.TypeOf(data))
	if err := json.Unmarshal(v, &data); err != nil {
		return nil, err
	}

	log.Printf("DefaultSerializer: Generated [ %v ]", data)
	return data, nil
}
