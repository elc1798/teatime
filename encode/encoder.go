package encode

// Write interface for serializer
type Serializer interface {
	Serialize(v interface{}) ([]byte, error)
	Deserialize([]byte) (interface{}, error)
}
