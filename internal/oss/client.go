package oss

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"

	"online-learning-platform/internal/config"
)

// Client 封装OSS客户端
type Client struct {
	client     *oss.Client
	bucket     *oss.Bucket
	bucketName string
}

var defaultClient *Client

// InitOSSClient 初始化OSS客户端
func InitOSSClient(cfg config.OSSConfig) error {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.BucketName == "" {
		return fmt.Errorf("invalid OSS config, please check endpoint/accessKey/bucket")
	}

	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return fmt.Errorf("failed to create OSS client: %w", err)
	}

	bucket, err := client.Bucket(cfg.BucketName)
	if err != nil {
		return fmt.Errorf("failed to get bucket %s: %w", cfg.BucketName, err)
	}

	defaultClient = &Client{
		client:     client,
		bucket:     bucket,
		bucketName: cfg.BucketName,
	}
	return nil
}

// GetClient 获取默认客户端
func GetClient() (*Client, error) {
	if defaultClient == nil {
		return nil, fmt.Errorf("OSS client not initialized")
	}
	return defaultClient, nil
}

// UploadReader 上传数据流
func UploadReader(ctx context.Context, objectKey string, reader io.Reader, options ...oss.Option) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	if err := client.bucket.PutObject(objectKey, reader, options...); err != nil {
		return "", fmt.Errorf("failed to upload object %s: %w", objectKey, err)
	}

	return fmt.Sprintf("https://%s.%s/%s", client.bucketName, client.client.Config.Endpoint, objectKey), nil
}

// UploadFile 上传本地文件
func UploadFile(ctx context.Context, objectKey, filePath string, options ...oss.Option) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	if err := client.bucket.PutObjectFromFile(objectKey, filePath, options...); err != nil {
		return "", fmt.Errorf("failed to upload file %s to object %s: %w", filePath, objectKey, err)
	}

	return fmt.Sprintf("https://%s.%s/%s", client.bucketName, client.client.Config.Endpoint, objectKey), nil
}

// GenerateSignedURL 生成签名URL
func GenerateSignedURL(objectKey string, expire time.Duration, httpMethod oss.HTTPMethod, options ...oss.Option) (string, error) {
	client, err := GetClient()
	if err != nil {
		return "", err
	}

	if expire <= 0 {
		expire = 15 * time.Minute
	}

	signedURL, err := client.bucket.SignURL(objectKey, httpMethod, int64(expire.Seconds()), options...)
	if err != nil {
		return "", fmt.Errorf("failed to sign url for object %s: %w", objectKey, err)
	}

	return signedURL, nil
}

// DeleteObject 删除对象
func DeleteObject(objectKey string) error {
	client, err := GetClient()
	if err != nil {
		return err
	}

	if err := client.bucket.DeleteObject(objectKey); err != nil {
		return fmt.Errorf("failed to delete object %s: %w", objectKey, err)
	}
	return nil
}
