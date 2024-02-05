package main

import (
	"Simp/src/rpc"
	"fmt"
)

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

// UserInfo
func (r *UserInfo) Decode(Bytes []byte) *UserInfo {
	d := new(rpc.Decode[UserInfo])
	d.ClassName = "UserInfo"
	d.Bytes = Bytes
	fmt.Println("Byte - ", Bytes)
	r.name = d.ReadString(1)
	fmt.Println("r.name - ", r.name)
	r.age = d.ReadInt8(2)
	fmt.Println("r.age - ", r.age)
	return r
}

func (r *UserInfo) Encode() *rpc.Encode[UserInfo] {
	d := new(rpc.Encode[UserInfo])
	d.ClassName = "UserInfo"
	d.Bytes = make([]byte, 1024)
	d.WriteString(1, r.name)
	d.WriteInt8(2, r.age)
	return d
}

// BasicInfo
func (r *BasicInfo) Decode(Bytes []byte) *BasicInfo {
	d := new(rpc.Decode[BasicInfo])
	d.ClassName = "BasicInfo"
	d.Bytes = Bytes
	r.Token = d.ReadString(1)
	return r
}

func (r *BasicInfo) Encode() *rpc.Encode[BasicInfo] {
	d := new(rpc.Encode[BasicInfo])
	d.ClassName = "BasicInfo"
	d.Bytes = make([]byte, 1024)
	d.WriteString(1, r.Token)
	return d
}

// User
func (r *User) Decode(Bytes []byte) *User {
	d := new(rpc.Decode[User])
	d.ClassName = "User"
	d.Bytes = Bytes
	r.BasicInfo = d.ReadStruct(1, "BasicInfo").(BasicInfo)
	r.UserInfo = d.ReadStruct(2, "UserInfo").(UserInfo)
	return r
}

func (r *User) Encode() *rpc.Encode[User] {
	d := new(rpc.Encode[User])
	d.ClassName = "User"
	d.Bytes = make([]byte, 1024)
	d.WriteStruct(1, r.BasicInfo)
	d.WriteStruct(2, r.UserInfo)
	return d
}

func main() {
	ui := new(UserInfo)
	ui.age = 1
	ui.name = "chelizichen"
	e := ui.Encode()
	ui.Decode(e.Bytes)
	//bi := new(BasicInfo)
	//bi.Token = "112213klsfnjiujas0218u321"
	//
	//u := new(User)
	//u.BasicInfo = *bi
	//u.UserInfo = *ui

}
