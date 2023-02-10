package localdb

import (
	"testing"

	"github.com/ccremer/clustercode/pkg/paperless"
	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	tests := map[string]struct {
		testFileName      string
		expectedDocuments map[int]paperless.Document
	}{
		"ExistingJSONFile": {
			testFileName: "test.metadata.json",
			expectedDocuments: map[int]paperless.Document{
				2:  {ID: 2},
				15: {ID: 15},
			},
		},
		"NonExistingJSONFile": {
			testFileName:      "nonexisting.metadata.json",
			expectedDocuments: map[int]paperless.Document{},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			old := fileName
			fileName = tt.testFileName
			defer func() {
				fileName = old
			}()

			result, err := Open("testdata")
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedDocuments, result.documents)
			assert.Equal(t, "testdata/"+tt.testFileName, result.filePath)
		})
	}
}

func TestDatabase_GetAll(t *testing.T) {
	tests := map[string]struct {
		givenDocuments    map[int]paperless.Document
		expectedDocuments []paperless.Document
	}{
		"NoDocuments": {
			givenDocuments:    map[int]paperless.Document{},
			expectedDocuments: []paperless.Document{},
		},
		"SeveralDocuments": {
			givenDocuments: map[int]paperless.Document{
				15: {ID: 15},
				1:  {ID: 1},
			},
			expectedDocuments: []paperless.Document{
				{ID: 1},
				{ID: 15},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			db := &Database{documents: tt.givenDocuments}
			result := db.GetAll()
			assert.Equal(t, tt.expectedDocuments, result)
		})
	}
}
