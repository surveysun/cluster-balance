package comm

import (
	"fmt"
	"testing"
)

func handle_put(key, value []byte) error {
	fmt.Println("[Put] key:", string(key), " value:", string(value))
	return nil
}

func handle_delete(key, value []byte) error {
	fmt.Println("[Delete] key:", string(key), " value:", string(value))
	return nil
}

func TestEtcdHander_Line(t *testing.T){
	etcd, err := NewEtcdHander([]string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Error("NewEtcdHander err:", err)
		t.Failed()
	}

	err = etcd.CreatePath("/home/sunwei/test-1")
	if err != nil {
		t.Error("CreatePath err:", err)
		t.Failed()
	}
	t.Log("Test CreatePath OK")

	kv, err := etcd.Delete("/home/sunwei/test-1")
	if err != nil {
		t.Error("Delete err:", err)
	}

	if kv.key != "/home/sunwei/test-1" {
		t.Error("the key not macth")
	}

	t.Log("Test Delete OK")
}

func TestEtcdHander_Watch(t *testing.T) {
	etcd, err := NewEtcdHander([]string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Error("NewEtcdHander err:", err)
	}

	go etcd.Watch("/home/sunwei/test-1", handle_put, handle_delete)

	etcd.Put("/home/sunwei/test-1", "test-1")
	//etcd.Delete("/home/sunwei/test-1")
}