// todo 2024.2.25 解码会出现乱码问题
package main

import (
	"Simp/src/rpc"
	"fmt"
	"reflect"
)

func structTest() {
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
	// u.Decode(ue.Bytes)
	//// 解码
	us := new(User)
	us.Decode(ue.Bytes)
	fmt.Println("main res", us.UserInfo.birth)
	fmt.Println("main res", us.UserInfo.age)
	fmt.Println("main res", us.UserInfo.name)
	fmt.Println("main res", us.BasicInfo.Token)
}

func baseArrTest() {
	{
		bi := new(BasicInfo)
		bi.Token = "lllll11111"
		qi := new(QueryIds)
		qi.BasicInfo = bi
		qi.Ids = []int32{1, 2, 3, 4, 5}
		e := qi.Encode()

		{
			fmt.Printf("e.Bytes: %v\n", e.Bytes)
			tt := new(QueryIds)
			tt.Decode(e.Bytes)
			fmt.Printf("tt.BasicInfo.Token: %v\n", tt.BasicInfo.Token)
			fmt.Printf("tt.Ids: %v\n", tt.Ids)
		}
	}
}

func structArrTest() {
	{
		ui := new(UserInfo)
		ui.age = 1
		ui.name = "chelizichenc"
		ui.birth = 32767

		usp := new(UserResp)
		usp.UserInfo = []UserInfo{*ui}
		e := usp.Encode()

		ur := new(UserResp)
		// fmt.Printf("e.Bytes: %v\n", e.Bytes)
		ur.Decode(e.Bytes)
		// fmt.Println("ur", ur.UserInfo[0])
	}
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
	r.BasicInfo = d.ReadStruct(1, reflect.ValueOf(new(BasicInfo))).Interface().(*BasicInfo)
	r.UserInfo = d.ReadStruct(2, reflect.ValueOf(new(UserInfo))).Interface().(*UserInfo)
	return r
}

func (r *User) Encode() *rpc.Encode[User] {
	d := new(rpc.Encode[User])
	d.ClassName = "User"
	d.Bytes = make([]byte, 1024)
	d.WriteStruct(1, reflect.ValueOf(r.BasicInfo))
	d.WriteStruct(2, reflect.ValueOf(r.UserInfo))
	return d
}

type QueryIds struct {
	BasicInfo *BasicInfo
	Ids       []int32
}

func (r *QueryIds) Decode(Bytes []byte) *QueryIds {
	d := new(rpc.Decode[QueryIds])
	d.ClassName = "QueryIds"
	d.Bytes = Bytes
	r.BasicInfo = d.ReadStruct(1, reflect.ValueOf(new(BasicInfo))).Interface().(*BasicInfo)
	r.Ids = d.ReadList(2, []int32{}).([]int32)
	return r
}

func (r *QueryIds) Encode() *rpc.Encode[QueryIds] {
	d := new(rpc.Encode[QueryIds])
	d.ClassName = "QueryIds"
	d.Bytes = make([]byte, 1024)
	d.WriteStruct(1, reflect.ValueOf(r.BasicInfo))
	d.WriteList(2, r.Ids)
	return d
}

type UserResp struct {
	UserInfo []UserInfo
}

func (r *UserResp) Decode(Bytes []byte) *UserResp {
	d := new(rpc.Decode[UserResp])
	d.ClassName = "UserResp"
	d.Bytes = Bytes
	r.UserInfo = d.ReadList(1, []UserInfo{}).([]UserInfo)
	return r
}

func (r *UserResp) Encode() *rpc.Encode[UserResp] {
	d := new(rpc.Encode[UserResp])
	d.ClassName = "UserResp"
	d.Bytes = make([]byte, 1024)
	d.WriteList(1, r.UserInfo)
	return d
}

type HashMapTest struct {
	HashTest map[string]interface{}
}

func (r *HashMapTest) Decode(Bytes []byte) *HashMapTest {
	d := new(rpc.Decode[HashMapTest])
	d.ClassName = "HashMapTest"
	d.Bytes = Bytes
	r.HashTest = d.ReadMap(1)
	return r
}

func (r *HashMapTest) Encode() *rpc.Encode[HashMapTest] {
	d := new(rpc.Encode[HashMapTest])
	d.ClassName = "HashMapTest"
	d.Bytes = make([]byte, 1024)
	d.WriteMap(1, r.HashTest)
	return d
}

func StructTest() {
	h := new(HashMapTest)
	h.HashTest = make(map[string]interface{})
	h.HashTest["greet"] = []int{1, 2, 3, 4, 5}
	h.HashTest["foo"] = []int{2, 3, 4, 5, 6, 2, 3, 4, 5}
	//h.HashTest["greet"] = "hello"
	//h.HashTest["foo"] = "bar"
	h.Encode()

}

func main() {
	// baseArrTest()
	// structArrTest()
	// structTest()
	StructTest()
}
