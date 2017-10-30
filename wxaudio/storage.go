package wxaudio

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
)

type Storage struct {
	refresh  time.Duration
	instance interface{}
}

type StorageUploader interface {
	Upload(string, []byte, map[string]string) (interface{}, error)
}

type StorageLinker interface {
	Link() error
}

type StorageUploadReturn struct {
	Hash string `json:"id"`
	Key  string `json:"filename"`
	Url  string `json:"url"`
}

func NewStorage() *Storage {
	return &Storage{
	//todo
	}
}

func (c *Storage) Set(instance interface{}) {
	c.instance = instance
}

func (c *Storage) Link() (err error) {
	err = c.instance.(StorageLinker).Link()
	return
}

func (c *Storage) Upload(filename string, data []byte, params map[string]string) (interface{}, error) {
	return c.instance.(StorageUploader).Upload(filename, data, params)
}

type QiniuInstance struct {
	accessKey string
	secretKey string
	bucket    string
	domain    string
	cfg       *storage.Config
}

func NewQiniuInstance(access, secret, bucket, domain string, cfg *storage.Config) *QiniuInstance {
	return &QiniuInstance{
		accessKey: access,
		secretKey: secret,
		bucket:    bucket,
		domain:    domain,
		cfg:       cfg,
	}
}

func (c *QiniuInstance) Link() error {
	return nil
}

func (c *QiniuInstance) upToken() string {
	putPolicy := storage.PutPolicy{
		Scope: c.bucket,
	}
	mac := qbox.NewMac(c.accessKey, c.secretKey)
	return putPolicy.UploadToken(mac)
}

func (c *QiniuInstance) Upload(filename string, data []byte, params map[string]string) (interface{}, error) {
	var res storage.PutRet
	formUploader := storage.NewFormUploader(c.cfg)
	err := formUploader.Put(context.Background(), &res, c.upToken(), filename, bytes.NewReader(data), int64(len(data)), &storage.PutExtra{
		Params: params,
	})
	var ret StorageUploadReturn
	if err == nil {
		ret.Hash = res.Hash
		ret.Key = res.Key
		ret.Url = fmt.Sprintf("https://%s/%s", c.domain, ret.Key)
	}
	return ret, err
}
