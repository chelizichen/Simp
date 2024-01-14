#!/bin/bash  

# if permission denied
# run script with ` chmod +x build.sh ` 

# compile
go build -o service_go

# build
tar -cvf CalcServer.tar.gz ./simp.yaml ./service_go