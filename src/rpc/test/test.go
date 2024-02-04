package main

import "Simp/src/rpc"

type UserInfo struct {
	age  int8
	name string
}

type BasicInfo struct {
	Token string
}

type User struct {
	UserInfo  UserInfo
	BasicInfo BasicInfo
}

func (r UserInfo) Decode() {
	d := new(rpc.Decode[UserInfo])
	d.ClassName = "UserInfo"
	d.Bytes = make([]byte, 1024)
	r.name = d.ReadString(1)
	r.age = d.ReadInt8(2)
}

func (r BasicInfo) Decode() {
	d := new(rpc.Decode[BasicInfo])
	d.ClassName = "User"
	d.Bytes = make([]byte, 1024)
	r.Token = d.ReadString(1)
}

func (r User) Decode() {
	d := new(rpc.Decode[User])
	d.ClassName = "User"
	d.Bytes = make([]byte, 1024)
	r.BasicInfo = d.ReadStruct(1, "BasicInfo").(BasicInfo)
	r.UserInfo = d.ReadStruct(2, "UserInfo").(UserInfo)
}
