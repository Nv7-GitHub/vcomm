package vcomm

import "reflect"

type VComm struct {
	server reflect.Value
}

func NewVComm(server interface{}) *VComm {
	return &VComm{
		server: reflect.ValueOf(server),
	}
}
