// Code generated by tools/cmd/genjwt/main.go. DO NOT EDIT.

package jwt

import (
	"fmt"
	"time"
)

// Builder is a convenience wrapper around the New() constructor
// and the Set() methods to assign values to Token claims.
// Users can successively call Claim() on the Builder, and have it
// construct the Token when Build() is called. This alleviates the
// need for the user to check for the return value of every single
// Set() method call.
// Note that each call to Claim() overwrites the value set from the
// previous call.
type Builder struct {
	claims []*ClaimPair
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Claim(name string, value interface{}) *Builder {
	b.claims = append(b.claims, &ClaimPair{Key: name, Value: value})
	return b
}

func (b *Builder) Audience(v []string) *Builder {
	return b.Claim(AudienceKey, v)
}

func (b *Builder) Expiration(v time.Time) *Builder {
	return b.Claim(ExpirationKey, v)
}

func (b *Builder) IssuedAt(v time.Time) *Builder {
	return b.Claim(IssuedAtKey, v)
}

func (b *Builder) Issuer(v string) *Builder {
	return b.Claim(IssuerKey, v)
}

func (b *Builder) JwtID(v string) *Builder {
	return b.Claim(JwtIDKey, v)
}

func (b *Builder) NotBefore(v time.Time) *Builder {
	return b.Claim(NotBeforeKey, v)
}

func (b *Builder) Subject(v string) *Builder {
	return b.Claim(SubjectKey, v)
}

// Build creates a new token based on the claims that the builder has received
// so far. If a claim cannot be set, then the method returns a nil Token with
// a en error as a second return value
func (b *Builder) Build() (Token, error) {
	tok := New()
	for _, claim := range b.claims {
		if err := tok.Set(claim.Key.(string), claim.Value); err != nil {
			return nil, fmt.Errorf(`failed to set claim %q: %w`, claim.Key.(string), err)
		}
	}
	return tok, nil
}
