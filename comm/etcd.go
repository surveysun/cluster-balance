package comm

import (
	"context"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"time"
)


type EtcdHander struct {
	hosts     []string
	client    *clientv3.Client
}

type KV struct {
	Key     string
	Value   string
}

const (
	EventTypePut = clientv3.EventTypePut
	EventTypeDelete = clientv3.EventTypeDelete
)

type CallbackHanderWatchPut func(key, value []byte)error
type CallbackHanderWatchDelete func(key, value []byte)error

func NewEtcdHander(hosts []string)(*EtcdHander, error){
	if hosts == nil {
		return nil, errors.New("hosts is nil")
	}

	cfg := clientv3.Config{
		Endpoints:   hosts,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdHander{
		client: client,
		hosts: hosts,
	}, nil
}

func (e *EtcdHander)CreatePath(path string) error{
	_, err := e.client.Put(context.Background(), path, "", clientv3.WithIgnoreValue())
	return err
}

func (e *EtcdHander)Put(key string, value string)error{
	_, err := e.client.Put(context.Background(), key, value)
	return err
}

// if the key is not exist, return nil,nil
func (e *EtcdHander)Get(key string)( *KV, error){
	re, err := e.client.Get(context.Background(), key)
	if err != nil{
		return nil, err
	}

	if len(re.Kvs) >= 1 {
		if string(re.Kvs[0].Key) != key{
			return nil, errors.New("key not match")
		}
		return &KV{
			Key:    string(re.Kvs[0].Key),
			Value:  string(re.Kvs[0].Value),
		}, nil
	}
	return nil, nil
}

func (e *EtcdHander)GetAll(key string)([]KV, error){
	re, err := e.client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil{
		return nil, err
	}

	out := make([]KV, 0)
	for _, kv := range re.Kvs {
		out = append(out, KV{string(kv.Key), string(kv.Value),})
	}

	return out, nil
}

func (e *EtcdHander)Delete(key string)( *KV, error){
	re, err := e.client.Delete(context.Background(), key)
	if err != nil{
		return nil, err
	}

	if len(re.PrevKvs) >= 1 {
		return &KV{
			Key: string(re.PrevKvs[0].Key),
			Value: string(re.PrevKvs[0].Value),
		},nil
	}

	return nil, nil
}

func (e *EtcdHander)DeleteAll(key string)( *KV, error){
	re, err := e.client.Delete(context.Background(), key, clientv3.WithPrefix())
	if err != nil{
		return nil, err
	}

	if len(re.PrevKvs) >= 1 {
		return &KV{
			Key: string(re.PrevKvs[0].Key),
			Value: string(re.PrevKvs[0].Value),
		},nil
	}

	return nil, nil
}

func (e *EtcdHander)Watch(key string) clientv3.WatchChan {
	return e.client.Watch(context.Background(), key)
}

func (e *EtcdHander)CloseWatchAll() error {
	return e.client.Watcher.Close()
}