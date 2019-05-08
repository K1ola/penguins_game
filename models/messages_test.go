package models

import (
	"testing"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
)

func TestString(t *testing.T) {
	user := &User{}
	user.Reset()
	user.String()
	user.ProtoMessage()
	user.Descriptor()
	user.XXX_Unmarshal([]byte(""))
	user.XXX_Marshal([]byte(""), true)
	user.XXX_Size()
	user.XXX_DiscardUnknown()
	var message proto.Message
	user.XXX_Merge(message)
	user.GetCount()
	user.GetEmail()
	user.GetHashPassword()
	user.GetID()
	user.GetPassword()
	user.GetLogin()
	user.GetPicture()
	user.GetScore()
	user = nil

	user.GetCount()
	user.GetEmail()
	user.GetHashPassword()
	user.GetID()
	user.GetPassword()
	user.GetLogin()
	user.GetPicture()
	user.GetScore()

	jwt := &JWT{}
	jwt.Reset()
	jwt.String()
	jwt.ProtoMessage()
	jwt.Descriptor()
	jwt.XXX_Unmarshal([]byte(""))
	jwt.XXX_Marshal([]byte(""), true)
	jwt.XXX_Size()
	jwt.XXX_DiscardUnknown()
	jwt.XXX_Merge(message)
	jwt.GetToken()
	jwt = nil
	jwt.GetToken()

	nothing := &Nothing{}
	nothing.Reset()
	nothing.String()
	nothing.ProtoMessage()
	nothing.Descriptor()
	nothing.XXX_Unmarshal([]byte(""))
	nothing.XXX_Marshal([]byte(""), true)
	nothing.XXX_Size()
	nothing.XXX_DiscardUnknown()
	nothing.XXX_Merge(message)

	var cc *grpc.ClientConn
	_ = NewAuthCheckerClient(cc)

}
