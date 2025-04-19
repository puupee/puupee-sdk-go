package puupeesdk

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/puupee/puupee-api-go"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

type ReleaseOp struct {
	api *puupee.APIClient
}

func NewReleaseOp(api *puupee.APIClient) *ReleaseOp {
	return &ReleaseOp{
		api: api,
	}
}

type CreateReleasePayload struct {
	AppName       string
	Version       string
	Notes         string
	Platform      string
	ArtifactType  string
	Channel       string
	Environment   string
	IsEnabled     bool
	IsForceUpdate bool
	// 二进制文件
	Filepath string
}

func FileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func SliceMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	buf := make([]byte, 256*1024)
	_, err = file.Read(buf)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	r := bytes.NewBuffer(buf)
	_, _ = io.Copy(hash, r)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

type RapidCode struct {
	Size     int64
	MD5      string
	SliceMD5 string
	Name     string
}

func NewRapidCodeFromFile(fpath string) *RapidCode {
	rc := &RapidCode{}
	fileMD5, err := FileMD5(fpath)
	if err != nil {
		panic(err)
	}
	sliceMD5, err := SliceMD5(fpath)
	if err != nil {
		panic(err)
	}
	info, err := os.Stat(fpath)
	if err != nil {
		panic(err)
	}
	rc.MD5 = fileMD5
	rc.Name = filepath.Base(fpath)
	rc.SliceMD5 = sliceMD5
	rc.Size = info.Size()
	return rc
}

func (rc *RapidCode) ID() string {
	return fmt.Sprintf("%s#%s#%d", rc.MD5, rc.SliceMD5, rc.Size)
}

func (rc *RapidCode) String() string {
	return fmt.Sprintf("%s#%s#%d#%s", rc.MD5, rc.SliceMD5, rc.Size, rc.Name)
}

func (rc *RapidCode) Key() string {
	return "files/" + base64.StdEncoding.EncodeToString([]byte(rc.ID())) + filepath.Ext(rc.Name)
}

func (op *ReleaseOp) Create(payload *CreateReleasePayload) error {
	appInfo, _, err := op.api.AppApi.GetByName(context.Background()).
		Name(payload.AppName).
		Execute()
	if err != nil {
		return err
	}
	platform := puupee.AppPlatform(payload.Platform)
	artifactType := puupee.ArtifactType(payload.ArtifactType)
	dto := puupee.CreateOrUpdateAppReleaseDto{
		Version:       &payload.Version,
		Notes:         &payload.Notes,
		Platform:      &platform,
		ArtifactType:  &artifactType,
		IsForceUpdate: &payload.IsForceUpdate,
		AppId:         appInfo.Id,
		IsEnabled:     &payload.IsEnabled,
		Channel:       &payload.Channel,
		Environment:   &payload.Environment,
	}
	if payload.Filepath != "" {
		fileMD5, err := FileMD5(payload.Filepath)
		if err != nil {
			panic(err)
		}
		sliceMD5, err := SliceMD5(payload.Filepath)
		if err != nil {
			panic(err)
		}
		info, err := os.Stat(payload.Filepath)
		if err != nil {
			panic(err)
		}
		filename := filepath.Base(payload.Filepath)
		filename = strings.ReplaceAll(filename, "+", ".")
		size := info.Size()
		dto.Md5 = &fileMD5
		dto.SliceMd5 = &sliceMD5
		dto.Size = &size
		key := fmt.Sprintf("apps/%s/%s/releases/%s/%s", *appInfo.Creator.UserName, *appInfo.Name, strings.ToLower(payload.Platform), filename)
		dto.Key = &key

		creditResult, _, err := op.api.AppApi.
			GetUploadCredentials(context.Background()).
			Key(key).
			Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `FileApi.ApiAppFilePreSignUrlPost``: %v\n", err)
			return err
		}
		if creditResult != nil {
			// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。 https://console.cloud.tencent.com/cos5/bucket
			// 替换为用户的 region，存储桶region可以在COS控制台“存储桶概览”查看 https://console.cloud.tencent.com/ ，关于地域的详情见 https://cloud.tencent.com/document/product/436/6224 。
			u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", *creditResult.BucketName, *creditResult.RegionId))
			b := &cos.BaseURL{BucketURL: u}
			c := cos.NewClient(b, &http.Client{
				Transport: &cos.AuthorizationTransport{
					// 通过环境变量获取密钥
					// 环境变量 COS_SECRETID 表示用户的 SecretId，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
					SecretID: creditResult.GetAccessKeyId(),
					// 环境变量 COS_SECRETKEY 表示用户的 SecretKey，登录访问管理控制台查看密钥，https://console.cloud.tencent.com/cam/capi
					SecretKey:    creditResult.GetAccessKeySecret(),
					SessionToken: creditResult.GetSecurityToken(),
					Expire:       time.Duration(creditResult.GetExpiredTime()) * time.Second,
					// Debug 模式，把对应 请求头部、请求内容、响应头部、响应内容 输出到标准输出
					Transport: &debug.DebugRequestTransport{
						RequestHeader: true,
						// Notice when put a large file and set need the request body, might happend out of memory error.
						RequestBody:    false,
						ResponseHeader: true,
						ResponseBody:   false,
					},
				},
			})

			opt := &cos.ObjectPutOptions{
				ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
					Listener: &cos.DefaultProgressListener{},
				},
			}
			uploadResult, err := c.Object.PutFromFile(context.Background(), key, payload.Filepath, opt)
			if err != nil {
				return err
			}
			bts, err := ioutil.ReadAll(uploadResult.Body)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", string(bts))
		}
	}

	resp, _, err := op.api.AppReleaseApi.
		CreateAppRelease(context.Background()).
		Body(dto).
		Execute()
	if err != nil {
		return err
	}
	PrintObject(resp)
	return nil
}

func (op *ReleaseOp) List(appName string) (*puupee.AppReleaseDtoPagedResultDto, error) {
	appDto, _, err := op.api.AppApi.GetByName(context.Background()).Name(appName).Execute()
	if err != nil {
		return nil, err
	}
	dto, _, err := op.api.AppReleaseApi.GetAppReleaseList(context.Background()).AppId(*appDto.Id).Execute()
	if err != nil {
		return nil, err
	}
	return dto, nil
}
