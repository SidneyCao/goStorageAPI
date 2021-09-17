# goStorageAPI  

storage golang API，用于监听rsync日志，并批量异步上传文件到对应bucket。  
目前仅支持GCP Storage，后续计划添加AWS S3。

# 启动方法：  
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
4. cd ../
5. 修改 getPara.sh 中的内容，添加对应的游戏名，rsync源
5. sh rsyncLogGCP.sh *{GAME_NAME}*

