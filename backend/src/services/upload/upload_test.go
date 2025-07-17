package upload_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ivolejon/pivo/services/document_loader"
	"github.com/ivolejon/pivo/services/upload"
	"github.com/stretchr/testify/require"
)

var (
	clientID  = uuid.MustParse("0ef8a743-f92b-4280-937b-ef1e4736c626")
	projectID = uuid.MustParse("b15377e4-60f1-11f0-9ce3-834692c66f23")
)

func TestUploadServiceNew(t *testing.T) {
	_, err := upload.NewUploadService(clientID, projectID)
	require.NoError(t, err)
}

func TestUploadServiceSave(t *testing.T) {
	svc, err := upload.NewUploadService(clientID, projectID)
	require.NoError(t, err)
	IDs, err := svc.Save(upload.UploadFileParams{Data: []byte("test"), Filename: "test.txt"})
	require.NoError(t, err)
	require.NotNil(t, IDs)
}

func TestUploadServiceSaveWrongFileType(t *testing.T) {
	svc, err := upload.NewUploadService(clientID, projectID)
	require.NoError(t, err)
	IDs, err := svc.Save(upload.UploadFileParams{Data: []byte("test"), Filename: "test.doccc"})
	require.ErrorIs(t, err, document_loader.ErrFileTypeNotSupported)
	require.Nil(t, IDs)
}
