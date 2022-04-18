package graphql

import (
	"encoding/json"
	"fmt"
	"strings"
)

type NearVectorArgumentBuilder struct {
	vector        []float32
	withCertainty bool
	certainty     float32
}

// WithVector sets the search vector to be used in query
func (b *NearVectorArgumentBuilder) WithVector(vector []float32) *NearVectorArgumentBuilder {
	b.vector = vector
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearVectorArgumentBuilder) WithCertainty(certainty float32) *NearVectorArgumentBuilder {
	b.withCertainty = true
	b.certainty = certainty
	return b
}

// Build build the given clause
func (b *NearVectorArgumentBuilder) build() string {
	clause := []string{}
	if b.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", b.certainty))
	}
	if b.vector != nil {
		vectorB, err := json.Marshal(b.vector)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal nearVector search vector: %s", err))
		}
		clause = append(clause, fmt.Sprintf("vector: %s", string(vectorB)))
	}
	return fmt.Sprintf("nearVector:{%v}", strings.Join(clause, " "))
}
