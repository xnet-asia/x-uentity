package cryptography

// Codec transforms partner request and response payloads.
type Codec interface {
	Encode(payload []byte) ([]byte, error)
	Decode(payload []byte) ([]byte, error)
}
