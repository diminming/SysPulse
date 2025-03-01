package client

import (
	"context"

	"github.com/syspulse/tracker/linux/common"
	"go.uber.org/zap"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var client *minio.Client

func InitFileServer() {
	var (
		endpoint        string = common.SysArgs.Storage.FileServer.Endpoint
		accessKeyID     string = common.SysArgs.Storage.FileServer.AccessKey
		secretAccessKey string = common.SysArgs.Storage.FileServer.SecretKey
		useSSL          bool   = common.SysArgs.Storage.FileServer.UseSSL
		bucketName      string = common.SysArgs.Storage.FileServer.BucketName
	)
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		zap.L().Fatal("error init minio client.", zap.Error(err))
	}
	client = minioClient
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		zap.L().Fatal("Failed to check bucket existence.", zap.Error(err))
	}

	if !exists {
		zap.L().Fatal("target bucket is not exists.", zap.String("bucket", bucketName))
	}
}

func CreateBucket(bucketName string) {

	ctx := context.Background()

	// Check to see if we already own this bucket (which happens if you run this twice)
	exists, errBucketExists := client.BucketExists(ctx, bucketName)
	if errBucketExists == nil && exists {
		zap.L().Info("We already own the bucket.", zap.String("bucket", bucketName))
	} else if errBucketExists != nil {
		zap.L().Fatal("Failed to check bucket existence.", zap.Error(errBucketExists))
	} else {
		err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			zap.L().Fatal("Failed to create bucket.", zap.Error(err))
		}
		zap.L().Info("Successfully created bucket.", zap.String("bucket", bucketName))
	}
}

func Upload2FileServer(bucketName, objectName, filePath, contentType string) {

	ctx := context.Background()
	// Upload the test file with FPutObject
	info, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		zap.L().Error("error upload to fileserver.", zap.Error(err))
	}

	zap.L().Info("Successfully uploaded result file.", zap.String("target name", objectName), zap.Int64("size", info.Size))
}
