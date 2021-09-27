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
	"strings"
	"sync"
	"time"

	storage "cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

//flag
var (
	method = flag.String("m", "list", "方法名\nlist 列出bucket下的所有objects\nupload 上传文件\n")
	bucket = flag.String("b", "", "bucket名 (默认为空)")
	files  = flag.String("f", "", "文件列表 (默认为空)")
	cache  = flag.String("c", "true", "是否缓存")
	prefix = flag.String("p", "", "需要移除的文件前缀 (默认为空)")
	gNum   = flag.Int("g", 5, "最大goroutine数量")
)

//缓存header内容
const (
	cacheMeta   string = "public, max-age=864000"
	noCacheMeta string = "no-store"
)

//创建wait group
var waitGroup sync.WaitGroup

//错误日志输出到stderr
var logerr *log.Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

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
func Upload(c *storage.Client, bucket string, jobChan chan string) {

	defer waitGroup.Done()
	for {

		file, ok := <-jobChan
		if !ok {
			break
		}
		object := strings.TrimPrefix(file, *prefix)
		//根据后缀检测Content-Type
		fileArray := strings.Split(file, ".")
		mtype := mime.TypeByExtension("." + fileArray[len(fileArray)-1])
		if mtype == "" {
			//默认 Content-Type
			mtype = "application/octet-stream"
		}
		//读取单个文件
		f, err := os.Open(file)
		if err != nil {
			logerr.Printf("failed to open %v: %v\n", file, err)
			return
		}
		defer f.Close()

		ctx := context.Background()
		//获取object
		o := c.Bucket(bucket).Object(object)
		//上传文件
		w := o.NewWriter(ctx)
		if _, err := io.Copy(w, f); err != nil {
			logerr.Printf("failed to uplaod %v: %v", file, err)
			return
		}
		if err := w.Close(); err != nil {
			logerr.Printf("Writer.Close: %v", err)
			return
		}

		//更新metadata
		objectAttrsToUpdate := storage.ObjectAttrsToUpdate{}
		objectAttrsToUpdate.ContentType = mtype
		objectAttrsToUpdate.CacheControl = cacheMeta
		if *cache == "false" {
			objectAttrsToUpdate.CacheControl = noCacheMeta
		}

		if _, err := o.Update(ctx, objectAttrsToUpdate); err != nil {
			logerr.Printf("failed to update metadata of %v: %v", object, err)
			return
		}

		log.Printf("successful to upload： %v\n", object)
	}
}

func main() {
	//获取命令行参数
	flag.Parse()

	//日志输出到stdout
	log.SetOutput(os.Stdout)

	//通过系统变量来进行认证
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/root/bucket-private.json")

	//新建client
	ctx := context.Background()
	c, err := storage.NewClient(ctx)
	if err != nil {
		logerr.Panicf("falied to create client: %v", err)
	}
	//关闭client
	defer c.Close()

	//wait group中始终有n+1个counter
	waitGroup.Add(1)

	//创建job队列
	jobChan := make(chan string, *gNum)

	switch *method {
	case "list":
		res, err := List(c, *bucket)
		if err != nil {
			logerr.Panicf("failed to list: %v", err)
		}
		for _, line := range res {
			fmt.Println(line)
		}
	case "upload":
		//启动*gNum个协程
		for i := 0; i < *gNum; i++ {
			go Upload(c, *bucket, jobChan)
		}
		//待上传文件以列表形式存储在文件中
		//按行读取文件，每个文件以goroutines形式上传
		f, err := os.Open(*files)
		if err != nil {
			logerr.Panicf("failed to open list file: %v", err)
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
			jobChan <- string(line)
		}
	}
	//decrease 最后一个counter
	waitGroup.Done()
	waitGroup.Wait()
	close(jobChan)
	log.Println("上传完成")
}
