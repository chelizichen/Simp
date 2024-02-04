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

func (r *UserInfo) Decode() {
	d := new(rpc.Decode[UserInfo])
	d.ClassName = "UserInfo"
	d.Bytes = make([]byte, 1024)
	r.name = d.ReadString(1)
	r.age = d.ReadInt8(2)
}

func (r *UserInfo) Encode() {
	d := new(rpc.Encode[UserInfo])
	d.ClassName = "UserInfo"
	d.Bytes = make([]byte, 1024)
	d.WriteString(1, r.name)
	d.WriteInt8(2, r.age)
}

func (r *BasicInfo) Decode() {
	d := new(rpc.Decode[BasicInfo])
	d.ClassName = "User"
	d.Bytes = make([]byte, 1024)
	r.Token = d.ReadString(1)
}

func (r *BasicInfo) Encode() {
	d := new(rpc.Encode[BasicInfo])
	d.ClassName = "UserInfo"
	d.Bytes = make([]byte, 1024)
	d.WriteString(1, r.Token)
}

func (r *User) Decode() {
	d := new(rpc.Decode[User])
	d.ClassName = "User"
	d.Bytes = make([]byte, 1024)
	r.BasicInfo = d.ReadStruct(1, "BasicInfo").(BasicInfo)
	r.UserInfo = d.ReadStruct(2, "UserInfo").(UserInfo)
}

func (r *User) Encode() {
	d := new(rpc.Encode[User])
	d.ClassName = "User"
	d.Bytes = make([]byte, 1024)
	d.WriteStruct(1, r.BasicInfo)
	d.WriteStruct(2, r.UserInfo)
}

func main() {
	ui := new(UserInfo)
	ui.age = 1
	ui.name = "chelizichen"

	bi := new(BasicInfo)
	bi.Token = "112213klsfnjiujas0218u321"

	u := new(User)
	u.BasicInfo = *bi
	u.UserInfo = *ui

	ui.Encode()
}
