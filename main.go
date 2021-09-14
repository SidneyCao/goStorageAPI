package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"sync"
	"time"

	storage "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var (
	method = flag.String("m", "list", "方法名\nlist 列出bucket下的所有objects\nupload 上传文件\n")
	bucket = flag.String("b", "", "bucket名(默认为空)")
	files  = flag.String("f", "", "文件列表(默认为空)")
)

func main() {
	//获取命令行参数
	flag.Parse()
	//通过系统变量来进行认证
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/root/bucket-private.json")

	//新建client
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		log.Panicf("falied to create client: %v", err)
	}
	//关闭client
	defer c.Close()

	//创建wait group
	//wait group中始终有n+1个counter
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	switch *method {
	case "list":
		res, err := List(c, *bucket)
		if err != nil {
			log.Panicf("failed to list: %v", err)
		}
		for _, line := range res {
			fmt.Println(line)
		}
	case "upload":
		//待上传文件以列表形式存储在文件中
		//按行读取文件，每个文件以goroutines形式上传
		f, err := os.Open(*files)
		if err != nil {
			log.Panicf("failed to open list file: %v", err)
		}
		//记得关闭文件
		defer f.Close()
		//按行读取文件
		br := bufio.NewReader(f)
		for {
			line, _, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			waitGroup.Add(1)
			go Upload(c, *bucket, string(line), string(line), &waitGroup)
		}
	}
	//decrease 最后一个counter
	waitGroup.Done()
	waitGroup.Wait()
	log.Println("上传完成")
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

//上传单个文件
func Upload(c *storage.Client, bucket string, file string, object string, waitGroup *sync.WaitGroup) {
	//读取单个文件
	f, err := os.Open(file)
	if err != nil {
		log.Printf("failed to open %v: %v\n", file, err)
	}

	mtype := mime.TypeByExtension(".txt")
	fmt.Println(mtype)
	ctx := context.Background()

	wc := c.Bucket(bucket).Object(object).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		log.Printf("failed to uplaod %v: %v", file, err)
		waitGroup.Done()
		return
	}
	if err := wc.Close(); err != nil {
		log.Printf("Writer.Close: %v", err)
		waitGroup.Done()
		return
	}
	log.Printf("successful to upload： %v: %v\n", object, mtype)
	waitGroup.Done()
}
