#!/bin/bash
#protoc --go_out=. *.proto

function build_dir(){
	if [ -n "$1" ];then
		protoc --go_out=./$1 $1/*.proto
	else
		protoc --go_out=. *.proto
	fi
}

#添加新增的编译目录即可
build_dir
build_dir message


# grpc
#protoc --go_out=plugins=grpc:. server.proto