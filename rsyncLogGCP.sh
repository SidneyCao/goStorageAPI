#!/bin/bash

#准备工作

game=$1

rsyncLog=/var/log/rsyncd/${game}.log
if [[ ! -f rsyncLog ]]; then 
    echo "请输入正确的游戏名"
    exit 1
fi


scriptDir=$(cd `dirname $0`; pwd)
logDir=/data/syncLog/${game}
taskListDir=/data/taskList/${game}

if [[ ! -d ${logDir}]]; then
        mkdir -p ${logDir}
fi

if [[ ! -d ${taskListDir}]]; then
    mkdir -p ${taskListDir}
fi


#监听日志
tail -f -n0 ${rsyncLog}| while read line; do
        echo ${line}
        #if echo "${Line}"; then
        #        FileName=$(echo "${Line}" | cut -d] -f2 | sed "s/^ *//")
        #        if [[ ${FileName} == receiving*  ]];then
        #                dateUpload=`date "+%Y-%m-%d %H:%M:%S"`
        #                echo "{\"startTime\":\"${dateUpload}\",\"status\":\"uploading\"}"  > ${monitorJson}
        #        elif [[ ${FileName} != rsync* && ${FileName} != sent* ]];then
        #                check_log
        #        fi
        #fi
done
