package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type Config struct {
	Minio struct {
		Endpoint  string `yaml:"endpoint"`
		AccessId  string `yaml:"access_id"`
		SecretKey string `yaml:"secret_key"`
		Region    string `yaml:"region"`
		Bucket    string `yaml:"bucket"`
	}
}

func main() {
	configArg := flag.String("config", "./minio.yaml", "指定配置文件")
	flag.Parse()
	config, err := ReadYamlConfig(*configArg)
	if err != nil {
		log.Fatalln(err)
	}
	// 解析参数
	var argArray []string
	for _, arg := range os.Args {
		if strings.Contains(arg, "-conf") || strings.Contains(arg, "yaml") {
			continue
		}
		argArray = append(argArray, arg)
	}
	imagePaths := argArray[1:]
	ctx := context.Background()
	endpoint := config.Minio.Endpoint
	accessID := config.Minio.AccessId
	secretKey := config.Minio.SecretKey
	region := config.Minio.Region
	bucket := config.Minio.Bucket

	if len(endpoint) == 0 || len(accessID) == 0 || len(secretKey) == 0 || len(region) == 0 || len(bucket) == 0 {
		log.Fatalln("未能获取到Minio客户端信息")
	}

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessID, secretKey, ""),
		Secure: true,
		Region: region,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("发现[%v]张图片,开始上传...", len(imagePaths))
	var minioImagePaths []string
	for index, filePath := range imagePaths {
		ext := path.Ext(filePath)
		fileName := uuid.New().String() + ext // 利用 uuid 作为文件名, 避免重复
		contentType, _ := GetFileContentTypeWithPath(filePath)
		// 上传
		info, err := minioClient.FPutObject(ctx, bucket, fileName, filePath, minio.PutObjectOptions{ContentType: contentType})
		log.Printf("上传 %v/%v", index+1, len(imagePaths))
		if err != nil {
			log.Fatalln(err)
		} else {
			minioImagePaths = append(minioImagePaths, fmt.Sprintf("https://%v/%v/%v\n", endpoint, info.Bucket, info.Key))
		}
	}
	log.Println("上传成功")
	for _, path := range minioImagePaths {
		fmt.Printf(path)
	}
}

// ReadYamlConfig 解析 yaml 配置文件
func ReadYamlConfig(path string) (*Config, error) {
	conf := &Config{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err = yaml.NewDecoder(f).Decode(conf)
		return conf, err
	}
}

// GetFileContentTypeWithPath 检测文件类型
func GetFileContentTypeWithPath(path string) (string, error) {
	file, err := os.Open(path)
	// 检测文件类型只需前512个字节就足够了
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}
	// 利用这个包的方法上传无论是否匹配到都有返回值, 如果没匹配到的时候默认返回 "application/octet-stream"
	return http.DetectContentType(buffer), nil
}
