#!/bin/bash  

# if permission denied
# run script with ` chmod +x build.sh ` 
readonly ServerName="CalcServer"

# rm
rm ./$ServerName.tar.gz ./service_go

# compile
go build -o service_go

# build
tar -cvf $ServerName.tar.gz ./simp.yaml ./service_go ./web