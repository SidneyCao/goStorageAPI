#!/bin/bash

case ${game} in
        test)
                srcDir=/data/test
                gstore=gs://wl-test
                gStatus=1
                nocacheStatus=1
                nocacheFile=zip.min.js
        ;;
	xiuzhen)
                srcDir=/data/xiuzhen.17996cdn.net
                gstore=gs://xiuzhen-cdn
                gStatus=1
                nocacheStatus=0
                nocacheFile=none
        ;;
	ysr)
                srcDir=/data/ysr.17996cdn.net
                gstore=gs://ysr-cdn
                gStatus=1
                nocacheStatus=0
                nocacheFile=none
        ;;
	yisu)
                srcDir=/data/yisu.17996cdn.net
                gstore=gs://yisu-cdn
                gStatus=1
                nocacheStatus=1
                nocacheFile='^.*server.cfg$|announcement.cfg$|updateinfo.xml$'
        ;;
        baoxiang)
                srcDir=/data/baoxiang.17996cdn.net
                gstore=gs://baoxiang-17996cdn-net
                gStatus=1
                nocacheStatus=0
                nocacheFile=none
#                nocacheFile=".*test.txt|abc$"
        ;;	
        *)
                echo "${game}  is not exist"
        ;;
esac
