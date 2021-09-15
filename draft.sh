upload_gs_log(){
        Gstore=`echo "${Gstore}"|sed 's/gs:\/\///g'`
        
        elif [[ nocacheStatus -eq 1 ]];then
                if [[ ${FileName} == *${nocacheFile} ]];then
                        gsutil -m -h 'Cache-Control: public, max-age=0' cp -r ${SrcDir}${FileName} gs://${Gstore}/${FileName}
                else
                        gsutil -m -h 'Cache-Control: public, max-age=864000' cp -r ${SrcDir}${FileName} gs://${Gstore}/${FileName}
                fi
        else
                gsutil -m -h 'Cache-Control: public, max-age=864000' cp -r ${SrcDir}${FileName} gs://${Gstore}/${FileName}
        fi
}



function check_log(){
        source ${ScriptDir}/rsyncgetdir.sh
        if [[ -f ${SrcDir}${FileName} ]];then
                if [[ GStatus -eq 1 ]];then
                        upload_gs_log
                        check_output
                fi
        else
                echo "${SrcDir}${FileName} 文件或目录格式不正确,正确格式如/test/1/ ,/test/1.txt"
        fi
}

