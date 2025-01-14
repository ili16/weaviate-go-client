package graphql

import (
	"io"
)

type NearImageArgumentBuilder struct {
	image        string
	imageReader  io.Reader
	hasCertainty bool
	certainty    float32
	hasDistance  bool
	distance     float32
}

// WithImage base64 encoded image
func (b *NearImageArgumentBuilder) WithImage(image string) *NearImageArgumentBuilder {
	b.image = image
	return b
}

// WithReader the image file
func (b *NearImageArgumentBuilder) WithReader(imageReader io.Reader) *NearImageArgumentBuilder {
	b.imageReader = imageReader
	return b
}

// WithCertainty that is minimally required for an object to be included in the result set
func (b *NearImageArgumentBuilder) WithCertainty(certainty float32) *NearImageArgumentBuilder {
	b.hasCertainty = true
	b.certainty = certainty
	return b
}

// WithDistance that is minimally required for an object to be included in the result set
func (b *NearImageArgumentBuilder) WithDistance(distance float32) *NearImageArgumentBuilder {
	b.hasDistance = true
	b.distance = distance
	return b
}

// Build build the given clause
func (b *NearImageArgumentBuilder) build() string {
	builder := &nearMediaArgumentBuilder{
		mediaName:  "nearImage",
		mediaField: "image",
		data:       b.image,
		dataReader: b.imageReader,
	}
	if b.hasCertainty {
		builder.withCertainty(b.certainty)
	}
	if b.hasDistance {
		builder.withDistance(b.distance)
	}
	return builder.build()
}
