package contextionary

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// ConceptGetter builder to get weaviate concpets
type ConceptGetter struct {
	connection *connection.Connection
	concept    string
}

// WithConcept that should be retrieved
func (cg *ConceptGetter) WithConcept(concept string) *ConceptGetter {
	cg.concept = concept
	return cg
}

// Do get the concept
func (cg *ConceptGetter) Do(ctx context.Context) (*models.C11yWordsResponse, error) {
	path := fmt.Sprintf("/c11y/concepts/%v", cg.concept)
	responseData, responseErr := cg.connection.RunREST(ctx, path, http.MethodGet, nil)
	err := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var concepts models.C11yWordsResponse
	parseErr := responseData.DecodeBodyIntoTarget(&concepts)
	if parseErr != nil {
		return nil, except.NewDerivedWeaviateClientError(parseErr)
	}
	return &concepts, nil
}
