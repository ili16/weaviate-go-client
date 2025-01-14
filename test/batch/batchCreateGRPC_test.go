package batch

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v4/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/grpc"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/testenv"
)

func TestBatchCreate_gRPC_integration(t *testing.T) {
	err := testenv.SetupLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to start weaviate")
	}

	port, _, _ := testsuit.GetPortAndAuthPw()
	cfg := weaviate.Config{
		Host:   fmt.Sprintf("localhost:%v", port),
		Scheme: "http",
		GrpcConfig: &grpc.Config{
			Host: "localhost:50051",
		},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		require.Nil(t, err)
	}

	t.Run("gRPC batch import", func(t *testing.T) {
		tests := []struct {
			name          string
			className     string
			properties    []map[string]interface{}
			withCrossRefs bool
		}{
			{
				name:       "all primitive properties",
				className:  "AllProperties",
				properties: testsuit.AllPropertiesDataAsMap(),
			},
			{
				name:       "all primitive properties with nested objects",
				className:  "AllPropertiesWithNested",
				properties: testsuit.AllPropertiesDataWithNestedObjectsAsMap(),
			},
			{
				name:       "all primitive properties with nested array objects",
				className:  "AllPropertiesWithNestedArray",
				properties: testsuit.AllPropertiesDataWithNestedArrayObjectsAsMap(),
			},
			{
				name:          "all primitive properties with cross references (single and multi ref types)",
				className:     "AllPropertiesWithCrossRefs",
				properties:    testsuit.AllPropertiesDataWithCrossReferencesAsMap(),
				withCrossRefs: true,
			},
			{
				name:          "all primitive properties with cross references (single and multi ref types) with nested and nested array objects",
				className:     "AllPropertiesWithCrossRefs",
				properties:    testsuit.AllPropertiesDataWithCrossReferencesWithNestedArrayObjectsAsMap(),
				withCrossRefs: true,
			},
		}
		for _, tt := range tests {
			className := tt.className
			objects := testsuit.AllPropertiesObjectsWithData(className, tt.properties)
			data := tt.properties

			err := client.Schema().AllDeleter().Do(context.Background())
			require.Nil(t, err)

			testsuit.AllPropertiesSchemaCreate(t, client, className, tt.withCrossRefs)

			batchResultSlice, batchErrSlice := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
			assert.Nil(t, batchErrSlice)
			assert.NotNil(t, batchResultSlice)
			assert.Equal(t, 3, len(batchResultSlice))

			for i := range objects {
				objs, err := client.Data().ObjectsGetter().
					WithID(objects[i].ID.String()).
					WithClassName(objects[i].Class).
					Do(context.Background())
				require.NoError(t, err)
				require.Len(t, objs, 1)
				obj := objs[0]
				assert.Equal(t, className, obj.Class)
				props, ok := obj.Properties.(map[string]interface{})
				require.True(t, ok)
				require.NotNil(t, props)
				properties := data[i]
				require.Equal(t, len(props), len(properties))
				for propName := range properties {
					assert.NotNil(t, props[propName])
				}
			}
		}
	})

	err = testenv.TearDownLocalWeaviate()
	if err != nil {
		require.Nil(t, err, "failed to tear down weaviate")
	}
}
