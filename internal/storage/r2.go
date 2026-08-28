package storage

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	appconfig "go-backend-boilerplate/internal/config"
)

type R2 struct {
	client        *s3.Client
	bucket        string
	publicBaseURL string
}

func NewR2(ctx context.Context, cfg *appconfig.Config) (*R2, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.R2AccessKeyID, cfg.R2SecretAccessKey, "",
		)),
	)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID)
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &R2{
		client:        client,
		bucket:        cfg.R2Bucket,
		publicBaseURL: strings.TrimRight(cfg.R2PublicBaseURL, "/"),
	}, nil
}

func (r *R2) Save(ctx context.Context, key string, body io.Reader, size int64, contentType string) (string, error) {
	in := &s3.PutObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   body,
	}
	if contentType != "" {
		in.ContentType = aws.String(contentType)
	}
	if size > 0 {
		in.ContentLength = aws.Int64(size)
	}
	if _, err := r.client.PutObject(ctx, in); err != nil {
		return "", err
	}
	return r.URL(key), nil
}

func (r *R2) Delete(ctx context.Context, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (r *R2) URL(key string) string {
	return r.publicBaseURL + "/" + strings.TrimLeft(key, "/")
}
