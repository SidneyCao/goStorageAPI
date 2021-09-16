#!/bin/bash

#准备工作

game=$1

rsyncLog=/var/log/rsyncd/${game}.log
if [[ ! -f $rsyncLog ]]; then 
    echo "请输入正确的游戏名"
    exit 1
fi


scriptDir=$(cd `dirname $0`; pwd)
taskLogDir=/data/taskLog/${game}
taskListDir=/data/taskList/${game}
goApiDir=${scriptDir}/gcpStorageAPI

if [[ ! -d ${taskLogDir} ]]; then
        mkdir -p ${taskLogDir}
fi

if [[ ! -d ${taskListDir} ]]; then
    mkdir -p ${taskListDir}
fi

source ${scriptDir}/getPara.sh

gs=`echo "${gstore}"|sed 's/gs:\/\///g'`

function touchTask(){
        if [[ ${nocacheStatus} -eq 1 ]];then
                touch ${taskListDir}/${dateUpload}-$taskID-noCache
        fi
        touch ${taskListDir}/${dateUpload}-$taskID-cache
}

function addTask(){
        if [[ -f ${srcDir}/${fileName} ]];then
                if [[ ${nocacheStatus} -eq 1 ]] && [[ ${srcDir}/${fileName} == *${nocacheFile} ]];then 
                        echo ${srcDir}/${fileName} >> ${taskListDir}/${dateUpload}-$taskID-noCache
                else 
                        echo ${srcDir}/${fileName} >> ${taskListDir}/${dateUpload}-$taskID-cache
                fi
        else    
                echo ''${srcDir}/${fileName}' 文件不存在'
        fi
}

#监听日志
tail -f -n0 ${rsyncLog}| while read line; do
        fileName=$(echo "${line}" | cut -d] -f2 | sed "s/^ *//")
        taskID=$(echo "${line}" | cut -d] -f1 | cut -d[ -f2)
        if [[ ${fileName} == 'receiving file list'  ]];then
                dateUpload=`date "+%Y_%m%d_%s"`
                touchTask
                echo '开始任务 '${taskID}'' >> ${taskLogDir}/${dateUpload}-$taskID.log
        elif [[ ${fileName} == sent* ]];then
                echo '开始上传需要缓存的文件'
                ${goApiDir}/gcpStorageAPI -b ${gs} -f ${taskListDir}/${dateUpload}-$taskID-cache -m upload >> ${taskLogDir}/${dateUpload}-$taskID.log
                echo '开始上传不需要缓存的文件'
                ${goApiDir}/gcpStorageAPI -b ${gs} -f ${taskListDir}/${dateUpload}-$taskID-noCache -m upload -c false >> ${taskLogDir}/${dateUpload}-$taskID.log
        else
                addTask
        fi
done
