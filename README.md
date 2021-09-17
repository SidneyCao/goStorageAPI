# goStorageAPI  

storage golang API，用于批量异步上传文件到bucket。  
目前仅支持gcp storage，后续计划添加aws s3。

# 使用方法：  
1. git clone https://github.com/SidneyCao/goStorageAPI.git  
2. cd goStorageAPI/gcpStorageAPI
3. go build 
4. cd ../ 
5. sh rsyncLogGCP.sh *{GAME_NAME}*
