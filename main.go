package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	storage "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var (
	bucket = flag.String("b", "", "bucket 默认为空")
	file   = flag.String("f", "", "文件名 默认为空")
)

func main() {
	//获取命令行参数
	flag.Parse()
	//通过系统变量来进行认证
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/simoncao/bucket-private.json")

	//新建client
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Panicf("falied to create client: %v", err)
	}
	defer c.Close()

	//fmt.Println(List(c, *bucket))

	f, err := os.Open(*file)
	if err != nil {
		log.Panicf("failed to open file: %v", err)
	}
	err = Upload(c, *bucket, f, *file)
	if err != nil {
		log.Panicf("failed to open file: %v", err)
	}
}

//列出bucket下的object
func List(c *storage.Client, bucket string) ([]string, error) {
	//利用context设定超时时间
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	//设定filter
	query := &storage.Query{Prefix: ""}
	bkt := c.Bucket(bucket)
	var names []string
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list : %w", err)
		}
		names = append(names, attrs.Name)
	}
	return names, nil
}

//上传文件
func Upload(c *storage.Client, bucket string, file *os.File, object string) error {
	ctx := context.Background()

	wc := c.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	log.Printf("成功上传文件： %v\n", object)
	return nil
}
