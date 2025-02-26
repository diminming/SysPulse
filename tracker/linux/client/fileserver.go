package client

import (
	"context"
	"log"

	"github.com/syspulse/tracker/linux/common"
	"go.uber.org/zap"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	ctx             context.Context = context.Background()
	endpoint        string          = common.SysArgs.Storage.FileServer.Endpoint
	accessKeyID     string          = common.SysArgs.Storage.FileServer.AccessKey
	secretAccessKey string          = common.SysArgs.Storage.FileServer.SecretKey
	useSSL          bool            = common.SysArgs.Storage.FileServer.UseSSL
	bucketName      string          = common.SysArgs.Storage.FileServer.BucketName
	client          *minio.Client
)

func init() {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Default().Fatalln(err)
	}
	client = minioClient
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalln("Failed to check bucket existence:", err)
	}

	if !exists {
		log.Fatalf("target bucket '%s' is not exists.", bucketName)
	}
}

func CreateBucket(bucketName string) {
	err := client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func Upload2FileServer(bucketName, objectName, filePath, contentType string) {
	// Upload the test file with FPutObject
	info, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		zap.L().Error("error upload to fileserver.", zap.Error(err))
	}

	zap.L().Info("Successfully uploaded result file.", zap.String("target name", objectName), zap.Int64("size", info.Size))
}
