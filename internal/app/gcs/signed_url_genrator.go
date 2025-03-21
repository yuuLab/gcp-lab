package gcs

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
)

const _defaultTTLMinutes = 15 * time.Minute

type SignedURLGenerator struct {
	bucketName string
}

// NewSignedURLGenerator creates SignedURLGenerator.
func NewSignedURLGenerator(bucketName string) *SignedURLGenerator {
	return &SignedURLGenerator{
		bucketName: bucketName,
	}
}

// Generate generates object signed URL with PUT method.
func (s *SignedURLGenerator) Generate(objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	// Signing a URL requires credentials authorized to sign a URL. You can pass
	// these in through SignedURLOptions with one of the following options:
	//    a. a Google service account private key, obtainable from the Google Developers Console
	//    b. a Google Access ID with iam.serviceAccounts.signBlob permissions
	//    c. a SignBytes function implementing custom signing.
	// In this example, none of these options are used, which means the SignedURL
	// function attempts to use the same authentication that was used to instantiate
	// the Storage client. This authentication must include a private key or have
	// iam.serviceAccounts.signBlob permissions.
	opts := &storage.SignedURLOptions{
		Scheme: storage.SigningSchemeV4,
		Method: "PUT",
		Headers: []string{
			"Content-Type:application/octet-stream",
		},
		Expires: time.Now().Add(s.getTTLMinutes()),
	}

	u, err := client.Bucket(s.bucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("Bucket(%q).SignedURL: %w", s.bucketName, err)
	}

	return u, nil
}

func (s *SignedURLGenerator) getTTLMinutes() time.Duration {
	expstr := os.Getenv("TTL_MINUTES")
	exp, err := strconv.Atoi(expstr)
	if err != nil {
		slog.Warn("Failed to parse TTL_MINUTES env var, using default value.")
		// default value.
		return _defaultTTLMinutes
	}

	return time.Duration(exp) * time.Minute
}
