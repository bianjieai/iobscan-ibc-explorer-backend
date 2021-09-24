#!/bin/bash
# todo 需要根据环境配置变量

function use_dev(){
  export NODE_ENV=development&&
  export DB_ADDR=192.168.150.40:27017&&
  export DB_USER=iris&&
  export DB_PASSWD=irispassword&&
  export DB_DATABASE=iobscan-iris&&
  export LCD_ADDR=http://192.168.150.40:1317&&
  env
}

function use_qa(){
  export NODE_ENV=development&&
  export DB_ADDR=192.168.150.33:37017&&
  export DB_USER=csrb&&
  export DB_PASSWD=csrbpassword&&
  export DB_DATABASE=sync2&&
  export LCD_ADDR=http://192.168.150.32:2317&&
  env
}

function use_prod(){
  export NODE_ENV=development&&
  export DB_ADDR=192.168.150.33:37017&&
  export DB_USER=csrb&&
  export DB_PASSWD=csrbpassword&&
  export DB_DATABASE=sync2&&
  export LCD_ADDR=http://192.168.150.32:2317&&
  env
}


if [ $1 == "d" ] ; then
  use_dev
elif [ $1 == "q" ] ; then
  use_qa
elif [ $1 == "p" ] ; then
  use_prod
fi
