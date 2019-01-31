package comm

import (
	"fmt"
	"strconv"
	"testing"
	"time"
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
		t.Fatal("NewEtcdHander err:", err)
	}

	err = etcd.Put("/home/sunwei/test-1", "test-1")
	if err != nil {
		t.Fatal("CreatePath err:", err)
	}
	t.Log("Test Put OK")

	_, err = etcd.Delete("/home/sunwei/test-1")
	if err != nil {
		t.Fatal("Delete err:", err)
	}
	t.Log("Test Delete OK")
}

func TestEtcdHander_Watch(t *testing.T) {
	etcd, err := NewEtcdHander([]string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Fatal("NewEtcdHander err:", err)
	}

	go etcd.Watch("/home/sunwei/test-2", handle_put, handle_delete)

	etcd.Put("/home/sunwei/test-2", "test-111")
	kv, err := etcd.Get("/home/sunwei/test-2")
	if err != nil{
		t.Fatal("Get err:", kv)
	}
	t.Log("Get key:", kv.key, " value:", kv.value)

	err = etcd.Put("/home/sunwei/test-2", "test-222")
	if err != nil {
		t.Fatal("Put err:", err)

	}

	_, err = etcd.Delete("/home/sunwei/test-2")
	if err != nil {
		t.Fatal("Delete err:", err)
		t.Failed()
	}
	t.Log("Delete OK")

	kv, err = etcd.Get("/home/sunwei/test-2")
	if err != nil{
		t.Fatal("Get err:", kv)
	}
	if kv == nil {
		t.Log("test Get null key pass:")
	}

	for i := 0; i < 10; i++ {
		etcd.Put("/home/sunwei/test-2/" + strconv.Itoa(i), "test-value-" + strconv.Itoa(i))
	}

	kvs, err := etcd.GetAll("/home/sunwei/test-2")
	t.Log("len kvs:", len(kvs))
	for _, kv := range kvs{
		t.Log("key:", kv.key, " value:",kv.value)
	}
	if len(kvs) != 10 {
		t.Fatal("GetAll error")
	}

	etcd.DeleteAll("/home/sunwei/test-2")
	kvs, err = etcd.GetAll("/home/sunwei/test-2")
	if len(kvs) != 0 {
		t.Fatal("test delete failed")
	}

	time.Sleep(time.Second)
	err = etcd.CloseWatchAll()
	if err != nil{
		t.Fatal("CloseWatchAll err:", err)
	}
	t.Log("close watch ok")
}
