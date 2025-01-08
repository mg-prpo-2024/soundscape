package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
)

type Storage interface {
	GeneratePresignedURL(fileName string, expiry time.Duration) (string, error)
	DeleteFile(fileName string) error
}

type storage struct {
	credential *azblob.SharedKeyCredential
	config     *Config
	serviceURL string // TODO: make this a URL type
	client     *azblob.Client
}

var _ Storage = (*storage)(nil)

func NewStorage(config *Config) *storage {
	cred, err := azblob.NewSharedKeyCredential(config.AzureAccountName, config.AzureAccountKey)
	if err != nil {
		panic(fmt.Errorf("failed to create shared key credential for azure blob storage: %w", err))
	}

	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net", config.AzureAccountName)
	if config.AppEnv == "local" {
		serviceURL = fmt.Sprintf("http://localhost:10000/%s", config.AzureAccountName)
	}

	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, &azblob.ClientOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create client for azure blob storage: %w", err))
	}
	return &storage{
		credential: cred,
		client:     client,
		config:     config,
		serviceURL: serviceURL,
	}
}

func (s *storage) GeneratePresignedURL(fileName string, expiry time.Duration) (string, error) {
	permissions := sas.BlobPermissions{
		Write:  true,
		Create: true,
	}

	protocol := sas.ProtocolHTTPS
	if s.config.AppEnv == "local" {
		protocol = sas.ProtocolHTTPSandHTTP
	}

	expiryTime := time.Now().Add(expiry)
	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      protocol,
		ExpiryTime:    expiryTime,
		Permissions:   permissions.String(),
		ContainerName: s.config.AzureContainerName,
		BlobName:      fileName,
	}.SignWithSharedKey(s.credential)

	if err != nil {
		return "", fmt.Errorf("failed to generate SAS token: %w", err)
	}

	blobURL := fmt.Sprintf("%s/%s/%s?%s",
		s.serviceURL,
		s.config.AzureContainerName,
		fileName,
		sasQueryParams.Encode(),
	)

	return blobURL, nil
}

func (s *storage) DeleteFile(fileName string) error {
	_, err := s.client.DeleteBlob(context.Background(), s.config.AzureContainerName, fileName, &azblob.DeleteBlobOptions{})
	return err
}
