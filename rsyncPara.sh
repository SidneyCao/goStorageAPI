#!/bin/bash

case ${game} in
        test)
                SrcDir=/data/test
                Gstore=gs://wl-test
                GStatus=1
                nocacheStatus=1
                nocacheFile=index.html

        *)
                echo "${game}  is not exist"
        ;;
esac