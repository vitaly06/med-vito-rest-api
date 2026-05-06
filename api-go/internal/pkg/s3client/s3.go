package s3client

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	xdraw "golang.org/x/image/draw"
)

const (
	// Лимит, к которому стремимся после сжатия.
	targetCompressedImageBytes = 4 * 1024 * 1024
	maxImageSidePixels         = 1920
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
	originalName, contentType, body = compressImageIfNeeded(originalName, contentType, body)

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

func compressImageIfNeeded(originalName, contentType string, body []byte) (string, string, []byte) {
	if len(body) == 0 {
		return originalName, contentType, body
	}

	// Быстрая фильтрация: пробуем только для jpeg/png.
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(originalName), "."))
	ct := strings.ToLower(strings.TrimSpace(contentType))
	if ext == "" && ct != "" {
		if strings.HasSuffix(ct, "jpeg") || strings.HasSuffix(ct, "jpg") {
			ext = "jpg"
		}
		if strings.HasSuffix(ct, "png") {
			ext = "png"
		}
	}
	if ext != "jpg" && ext != "jpeg" && ext != "png" {
		return originalName, contentType, body
	}

	img, format, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		return originalName, contentType, body
	}
	format = strings.ToLower(format)
	if format != "jpeg" && format != "png" {
		return originalName, contentType, body
	}
	img = resizeImageIfNeeded(img, maxImageSidePixels)

	bestBody := body
	bestName := originalName
	bestCT := contentType

	// Для JPEG всегда пробуем лоссy-сжатие.
	qualities := []int{82, 75, 68, 60}
	for _, q := range qualities {
		var b bytes.Buffer
		if err := jpeg.Encode(&b, img, &jpeg.Options{Quality: q}); err != nil {
			continue
		}
		out := b.Bytes()
		if len(out) < len(bestBody) {
			bestBody = out
			bestName = replaceExtWithJPG(originalName)
			bestCT = "image/jpeg"
		}
		if len(bestBody) <= targetCompressedImageBytes {
			break
		}
	}

	// Для PNG с прозрачностью безопаснее остаться в PNG, но всё равно переупаковать.
	if hasAlphaChannel(img) {
		var b bytes.Buffer
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		if err := enc.Encode(&b, img); err == nil {
			out := b.Bytes()
			if len(out) < len(bestBody) {
				bestBody = out
				bestName = replaceExtWithPNG(originalName)
				bestCT = "image/png"
			}
		}
	}

	return bestName, bestCT, bestBody
}

func replaceExtWithJPG(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[:i] + ".jpg"
	}
	return name + ".jpg"
}

func replaceExtWithPNG(name string) string {
	if i := strings.LastIndex(name, "."); i >= 0 {
		return name[:i] + ".png"
	}
	return name + ".png"
}

func hasAlphaChannel(img image.Image) bool {
	b := img.Bounds()
	// Выборка по сетке, чтобы не проходить каждый пиксель на больших картинках.
	stepX := maxInt((b.Dx()/48)+1, 1)
	stepY := maxInt((b.Dy()/48)+1, 1)
	for y := b.Min.Y; y < b.Max.Y; y += stepY {
		for x := b.Min.X; x < b.Max.X; x += stepX {
			_, _, _, a := img.At(x, y).RGBA()
			if a != 0xffff {
				return true
			}
		}
	}
	return false
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func resizeImageIfNeeded(img image.Image, maxSide int) image.Image {
	if img == nil || maxSide <= 0 {
		return img
	}
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return img
	}
	if w <= maxSide && h <= maxSide {
		return img
	}

	newW, newH := w, h
	if w >= h {
		newW = maxSide
		newH = int(float64(h) * (float64(maxSide) / float64(w)))
	} else {
		newH = maxSide
		newW = int(float64(w) * (float64(maxSide) / float64(h)))
	}
	if newW < 1 {
		newW = 1
	}
	if newH < 1 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), img, b, xdraw.Over, nil)
	return dst
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
