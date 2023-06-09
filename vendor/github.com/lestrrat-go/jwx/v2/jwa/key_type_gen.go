// Code generated by tools/cmd/genjwa/main.go. DO NOT EDIT.

package jwa

import (
	"fmt"
	"sort"
	"sync"
)

// KeyType represents the key type ("kty") that are supported
type KeyType string

// Supported values for KeyType
const (
	EC             KeyType = "EC"  // Elliptic Curve
	InvalidKeyType KeyType = ""    // Invalid KeyType
	OKP            KeyType = "OKP" // Octet string key pairs
	OctetSeq       KeyType = "oct" // Octet sequence (used to represent symmetric keys)
	RSA            KeyType = "RSA" // RSA
)

var allKeyTypes = map[KeyType]struct{}{
	EC:       {},
	OKP:      {},
	OctetSeq: {},
	RSA:      {},
}

var listKeyTypeOnce sync.Once
var listKeyType []KeyType

// KeyTypes returns a list of all available values for KeyType
func KeyTypes() []KeyType {
	listKeyTypeOnce.Do(func() {
		listKeyType = make([]KeyType, 0, len(allKeyTypes))
		for v := range allKeyTypes {
			listKeyType = append(listKeyType, v)
		}
		sort.Slice(listKeyType, func(i, j int) bool {
			return string(listKeyType[i]) < string(listKeyType[j])
		})
	})
	return listKeyType
}

// Accept is used when conversion from values given by
// outside sources (such as JSON payloads) is required
func (v *KeyType) Accept(value interface{}) error {
	var tmp KeyType
	if x, ok := value.(KeyType); ok {
		tmp = x
	} else {
		var s string
		switch x := value.(type) {
		case fmt.Stringer:
			s = x.String()
		case string:
			s = x
		default:
			return fmt.Errorf(`invalid type for jwa.KeyType: %T`, value)
		}
		tmp = KeyType(s)
	}
	if _, ok := allKeyTypes[tmp]; !ok {
		return fmt.Errorf(`invalid jwa.KeyType value`)
	}

	*v = tmp
	return nil
}

// String returns the string representation of a KeyType
func (v KeyType) String() string {
	return string(v)
}
