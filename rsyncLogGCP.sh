#!/bin/bash

#准备工作

game=$1

rsyncLog=/var/log/rsyncd/${game}.log
if [[ ! -f $rsyncLog ]]; then 
    echo "请输入正确的游戏名"
    exit 1
fi


scriptDir=$(cd `dirname $0`; pwd)
logDir=/data/syncLog/${game}
taskListDir=/data/taskList/${game}

if [[ ! -d ${logDir} ]]; then
        mkdir -p ${logDir}
fi

if [[ ! -d ${taskListDir} ]]; then
    mkdir -p ${taskListDir}
fi

source ${scriptDir}/rsyncPara.sh

function touchTask(){
        if [[ ${nocacheStatus} -eq 1 ]];then
                touch ${taskListDir}/${dateUpload}-$taskID-noCache
        fi
        touch ${taskListDir}/${dateUpload}-$taskID-cache
}

function addTask(){
        if [[ ${nocacheStatus} -eq 1 ]] && [[ ${srcDir}/${fileName} == *${nocacheFile} ]];then 
                echo ${srcDir}/${fileName} >> ${taskListDir}/${dateUpload}-$taskID-noCache
        else 
                echo ${srcDir}/${fileName} >> ${taskListDir}/${dateUpload}-$taskID-cache
        fi
}

#监听日志
tail -f -n0 ${rsyncLog}| while read line; do
        fileName=$(echo "${line}" | cut -d] -f2 | sed "s/^ *//")
        taskID=$(echo "${line}" | cut -d] -f1 | cut -d[ -f2)
        if [[ ${fileName} == 'receiving file list'  ]];then
                dateUpload=`date "+%Y-%m-%d %H:%M:%S"`
                touchTask
        elif [[ ${fileName} != sent* ]];then
                echo "执行"
        else
                addTask
        fi
done
