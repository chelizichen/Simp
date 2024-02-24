// todo 2024.2.25 解码会出现乱码问题
package main

import (
	"Simp/src/rpc"
	"fmt"
)

func main() {
	ui := new(UserInfo)
	ui.age = 1
	ui.name = "chelizichenc"
	ui.birth = 32767
	e := ui.Encode()
	//
	//n_ui := new(UserInfo)
	//n_ui.Decode(e.Bytes)
	//fmt.Println(n_ui.birth)
	//fmt.Println(n_ui.age)
	//fmt.Println(n_ui.name)

	ui.Decode(e.Bytes)
	bi := new(BasicInfo)
	bi.Token = "112213klsfnjiujas0218u321"

	u := new(User)
	u.BasicInfo = bi
	u.UserInfo = ui
	ue := u.Encode()
	u.Decode(ue.Bytes)
	//// 解码
	us := new(User)
	us.Decode(ue.Bytes)
	fmt.Println("main res", us.UserInfo.birth)
	fmt.Println("main res", us.UserInfo.age)
	fmt.Println("main res", us.UserInfo.name)
	fmt.Println("main res", us.BasicInfo.Token)
}

type BasicInfo struct {
	Token string
}

type UserInfo struct {
	age   int8
	birth int16
	name  string
}

type User struct {
	UserInfo  *UserInfo
	BasicInfo *BasicInfo
}

// UserInfo
func (r *UserInfo) Decode(Bytes []byte) *UserInfo {
	d := new(rpc.Decode[UserInfo])
	d.ClassName = "UserInfo"
	d.Bytes = Bytes
	r.name = d.ReadString(1)
	r.age = d.ReadInt8(2)
	r.birth = d.ReadInt16(3)
	return r
}

func (r *UserInfo) Encode() *rpc.Encode[UserInfo] {
	d := new(rpc.Encode[UserInfo])
	d.ClassName = "UserInfo"
	d.Bytes = make([]byte, 1024)
	d.WriteString(1, r.name)
	d.WriteInt8(2, r.age)
	d.WriteInt16(3, r.birth)
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
	_basicInfo := new(BasicInfo)
	d.ReadStruct(1, _basicInfo)
	r.BasicInfo = _basicInfo
	_userInfo := new(UserInfo)
	d.ReadStruct(2, _userInfo)
	r.UserInfo = _userInfo
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
