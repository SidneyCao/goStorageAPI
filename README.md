# goStorageAPI  
自用小工具  
storage golang API，用于监听rsync日志，并批量异步上传文件到对应bucket。  
目前仅支持GCP Storage，后续计划添加AWS S3。

## 原理
监听rsync日志中的事件, 此处 taskID=13836  
```
2021/09/15 07:57:23 [13836] rsync to test from x@UNKNOWN (x.x.x.x)
2021/09/15 07:57:23 [13836] receiving file list
2021/09/15 07:57:23 [13836] test/sync.txt
...
2021/09/15 07:57:24 [13836] sent 2392 bytes  received 1227017 bytes  total size 1220788
```
判断文件是否需要缓存，
并将文件写入  
/data/taskList/*{GAME_NAME}*/*{%Y_%M%d_%s}*-*{taskID}*-cache  
或 /data/taskList/*{GAME_NAME}*/*{%Y_%M%d_%s}*-*{taskID}*-noCache  
调用gcp storage golang api 进行一次认证，并通过goroutine批量异步上传  


## 启动前准备：
1. 首先部署好rsyncd服务，rsync日志需要在 /var/log/rsyncd/*{GAME_NAME}*.log  
2. 下载对应key，放入对应的目录中  
```
    os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/root/bucket-private.json")
```

## 启动方法：  
1. git clone https://github.com/SidneyCao/goStorageAPI.git  
2. cd goStorageAPI/gcpStorageAPI
3. go build  
4. 编译完成后可以查看一下二进制的使用方法
```
    ./gcpStorageAPI -h
    Usage of ./gcpStorageAPI:
    -b string
    	bucket名 (默认为空)
    -c string
    	是否缓存 (default "true")
    -f string
    	文件列表 (默认为空)
    -m string
    	方法名
    	list 列出bucket下的所有objects
    	upload 上传文件
    	 (default "list")
    -p string
    	需要移除的文件前缀 (默认为空)
    -t int
    	最大协程数 (默认为5) (default 5)
```
5. 回到上层目录，cd ../
6. 修改 getPara.sh 中的内容，添加对应的游戏名，rsync源，bucket名，缓存选项
7. 开始监听rsync日志，sh rsyncLogGCP.sh *{GAME_NAME}*
8. 成功上传的日志会输入到stdout，错误日志会通过stderr输入到 /data/taskLog/*{GAME_NAME}*/下
9. 任务状态会实时更新到 /data/taskLog/*{GAME_NAME}*/result.html 中，后续可以通过nginx提供给合作方

## 注意事项  

rsync -av *{YOUR_FILE}* --port=*{PORT}* --password-file=secret *{USER}*@*{HOSTNAME}*::*{GAME_NAME}* 
上传后文件会自动同步到bucket  
测试链接：  
https://*{DOMAIN}*/*{YOUR_FILE}*  

请勿使用rsync本地的相对位置，这样上传的文件位置会不对  
举例：  
rsync -av test.txt --port=*{PORT}* --password-file=secret *{USER}*@*{HOSTNAME}*::*{GAME_NAME}  
上传后的链接：  
https://*{DOMAIN}*/test.txt  

rsync -av test/test.txt --port=*{PORT}* --password-file=secret *{USER}*@*{HOSTNAME}*::*{GAME_NAME}  
上传后的链接仍然会是：  
https://*{DOMAIN}*/test.txt  


正确的做法应该是：  
rsync -av test/test.txt --port=*{PORT}* --password-file=secret *{USER}*@*{HOSTNAME}*::*{GAME_NAME}/test/   
上传后的链接：  
https://*{DOMAIN}*/test/test.txt  


上传目录同理