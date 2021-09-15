#!/bin/bash

case ${game} in
        test)
                srcDir=/data/test
                gstore=gs://wl-test
                gStatus=1
                nocacheStatus=1
                nocacheFile=index.html

        *)
                echo "${game}  is not exist"
        ;;
esac