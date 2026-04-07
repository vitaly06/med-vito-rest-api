package s3client

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Client — S3-совместимое хранилище (Beget и т.п.), как Nest S3Service.
type Client struct {
	bucket     string
	publicBase string
	svc        *s3.Client
}

// New создаёт клиент; publicBase — префикс публичного URL (как в Nest).
func New(endpoint, region, bucket, accessKey, secretKey, publicBase string) (*Client, error) {
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("нужны S3_ACCESS_KEY и S3_SECRET_KEY")
	}
	if bucket == "" {
		return nil, fmt.Errorf("нужен S3_BUCKET_NAME")
	}
	if endpoint == "" {
		endpoint = "https://s3.ru1.storage.beget.cloud"
	}
	if region == "" {
		region = "ru1"
	}
	if publicBase == "" {
		publicBase = "https://s3.ru1.storage.beget.cloud"
	}
	publicBase = strings.TrimRight(publicBase, "/")

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		awsconfig.WithRequestChecksumCalculation(aws.RequestChecksumCalculationWhenRequired),
		awsconfig.WithResponseChecksumValidation(aws.ResponseChecksumValidationWhenRequired),
	)
	if err != nil {
		return nil, err
	}
	svc := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &Client{bucket: bucket, publicBase: publicBase, svc: svc}, nil
}

// Upload загружает байты в папку (например users/) и возвращает публичный URL как в Nest.
func (c *Client) Upload(ctx context.Context, folder, originalName, contentType string, body []byte) (string, error) {
	ext := "bin"
	if i := strings.LastIndex(originalName, "."); i >= 0 && i < len(originalName)-1 {
		ext = strings.ToLower(strings.TrimSpace(originalName[i+1:]))
	}
	key := fmt.Sprintf("%s/%s.%s", strings.Trim(folder, "/"), uuid.NewString(), ext)
	ct := contentType
	if ct == "" {
		ct = "application/octet-stream"
	}
	_, err := c.svc.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(body),
		ContentType:   aws.String(ct),
		ContentLength: aws.Int64(int64(len(body))),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/%s", c.publicBase, c.bucket, key), nil
}

// DeleteByURL удаляет объект по полному URL или по ключу (без http).
func (c *Client) DeleteByURL(ctx context.Context, fileURL string) error {
	key, err := extractObjectKey(fileURL, c.bucket)
	if err != nil || key == "" {
		return nil
	}
	_, err = c.svc.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	return err
}

func extractObjectKey(fileURL, bucket string) (string, error) {
	if !strings.Contains(fileURL, "http") {
		k := strings.TrimPrefix(fileURL, "/")
		if strings.HasPrefix(k, bucket+"/") {
			k = strings.TrimPrefix(k, bucket+"/")
		}
		return k, nil
	}
	parts := strings.SplitN(fileURL, ".cloud/", 2)
	var pathPart string
	if len(parts) > 1 {
		pathPart = parts[1]
	} else {
		u, err := url.Parse(fileURL)
		if err != nil {
			return "", err
		}
		pathPart = strings.TrimPrefix(u.Path, "/")
	}
	// Убираем префикс bucket/ из ключа для API PutObject/DeleteObject.
	if strings.HasPrefix(pathPart, bucket+"/") {
		pathPart = strings.TrimPrefix(pathPart, bucket+"/")
	}
	return pathPart, nil
}
