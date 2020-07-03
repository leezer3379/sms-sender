package corp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/toolkits/pkg/logger"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const dingTimeOut = time.Second * 100


// Client
type Client struct {
	Mobiles []string
	Message string
	openUrl string
}

// New
func New(mobiles []string, message string, openurl string) *Client {
	c := new(Client)

	c.openUrl = openurl
	c.Mobiles = mobiles
	c.Message = message
	return c
}

func (c *Client) GetMobiles() []string {
	return c.Mobiles
}

// Send 发送信息
func (c *Client) Send(mobile string, msg string) error {

	postData := c.generateData(mobile, msg)

	url := c.openUrl

	resultByte, err := jsonPost(url, postData)
	if err != nil {
		return fmt.Errorf("invoke send api fail: %v", err)
	}

	if string(resultByte) != "ok" {
		err = fmt.Errorf("invoke send api return = %s", resultByte)
	}

	return err
}

func jsonPost(url string, data interface{}) ([]byte, error) {
	jsonBody, err := encodeJSON(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		logger.Info("ding talk new post request err =>", err)
		return nil, err
	}

	//req.Header.Set("Content-Type", "application/json")

	client := getClient()
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("ding talk post request err =>", err)
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func encodeJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Client) generateData(mobile string, msg string) interface{} {
	postData := make(map[string]interface{})
	postData["mobile"] = mobile
	postData["message"] = msg
	return postData
}

func getClient() *http.Client {
	// 通过http.Client 中的 DialContext 可以设置连接超时和数据接受超时 （也可以使用Dial, 不推荐）
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				conn, err := net.DialTimeout(network, addr, dingTimeOut) // 设置建立链接超时
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(dingTimeOut)) // 设置接受数据超时时间
				return conn, nil
			},
			ResponseHeaderTimeout: dingTimeOut, // 设置服务器响应超时时间
		},
	}
}
