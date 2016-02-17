#!/bin/bash

#function that prints out usage syntax
syntax () {
    echo " "
    echo "./app-xxx.sh [options]"
    echo " "
    echo "  [options] is one of create | backup | rm | read "
    echo "    start             - start concurrent application "
    echo " "
    echo "    stop             - stop concurrent application"
    echo " "
    echo "    restart                 - restart concurren application"
    echo " "
    echo "    bstart               - build and start concurrent application"
    echo " "
    echo "    brestart               - build and restart concurrent application"
    echo " "
    echo " "
}              



port=''$2
runmode=''$3
application='app-'${runmode}'-'${port}

if echo "$port" | grep -q '^[0-9]\+$'; then
    echo "["$1"] application: $application ,port: $port."
else
    echo "$0 {run:start|stop|restart} {port:number} {runmode:dev|test|prod}"
    exit 4
fi

if [ $runmode != "dev" ] && [ $runmode != "test" ] && [ $runmode != "prod" ]; then
    echo "$0 {run:start|stop|restart} {port:number} {runmode:dev|test|prod}"
    exit 4
fi

function exct(){
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./${application}/${application} ./../
    cp -r ./../conf ./../static ./../views ./${application}/
    mkdir ./${application}/log
    sed -i '' "s/^httpport.*/httpport = ${port}/g" ./${application}/conf/app.conf
    sed -i '' "s/^runmode.*/runmode = ${runmode}/g" ./${application}/conf/app.conf
}


case $1 in
    start)
        nohup ./${application}/${application} 2>&1 >> ./${application}/run.log 2>&1 /dev/null &
        echo ${application}" 服务已启动..."
        sleep 1
    ;;
    restart)
        killall ${application}
        sleep 1
        nohup ./${application}/${application} 2>&1 >> ./${application}/run.log 2>&1 /dev/null &
        echo ${application}" 服务已重启..."
        sleep 1
    ;;
    bstart)
        exct
        echo ${application}" 服务构建..."
        sleep 1
        nohup ./${application}/${application} 2>&1 >> ./${application}/run.log 2>&1 /dev/null &
        echo ${application}" 服务已启动..."
        sleep 1
    ;;
    stop)
        killall ${application}
        echo ${application}" 服务已停止..."
        sleep 1
    ;;
    brestart)
        rm -rf ./${application}/
        killall ${application}
        sleep 1
        exct
        echo ${application}" 服务构建..."
        sleep 1
        nohup ./${application}/${application} 2>&1 >> ./${application}/run.log 2>&1 /dev/null &
        echo ${application}" 服务已重启..."
        sleep 1
    ;;
    *)
        echo "$0 {start|stop|restart}"
        exit 4
    ;;
esac


