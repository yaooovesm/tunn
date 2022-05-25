package authenticationv2

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	log "github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
	"tunn/config"
	"tunn/transmitter"
	"tunn/utils/timer"
)

/**
  离线场景：客户端意外断开(客户端断网)，服务器意外断开(服务器断网)，客户端主动断开，服务端主动断开
*/

//
// Client
// @Description:
//
type Client struct {
	UUID      string              //客户端唯一ID
	handler   AuthClientHandler   //客户端处理
	PublicKey []byte              //服务端公钥 (登录后接收)
	WSKey     string              //ws&wss接入点 (登录后接收)
	version   transmitter.Version //传输版本
	tunnel    *transmitter.Tunnel //客户端连接
	Online    bool
}

//
// NewClient
// @Description:
// @param handler
// @return client
// @return err
//
func NewClient(handler AuthClientHandler) (client *Client, err error) {
	//检查证书
	if config.Current.Security.CertPem == "" {
		return nil, ErrCertFileNotFound
	}
	//设置连接地址
	address := strings.Join([]string{config.Current.Auth.Address, strconv.Itoa(config.Current.Auth.Port)}, ":")
	log.Info("connect to authentication server : ", address)
	//读取证书
	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.Current.Security.CertPem)
	if err != nil {
		return nil, log.Error("cannot load cert : ", err)
	}
	pool.AppendCertsFromPEM(ca)
	//准备连接
	dialer := websocket.Dialer{
		TLSClientConfig:   &tls.Config{RootCAs: pool},
		HandshakeTimeout:  time.Second * time.Duration(45),
		EnableCompression: false,
	}
	u := url.URL{Scheme: "wss", Host: address, Path: "/authentication"}
	//连接到服务器
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		//TODO 连接失败,清空当前客户端
		_ = conn.Close()
		return nil, ErrConnectFailed
	}
	wsConn := transmitter.WrapWSConn(conn)
	v := transmitter.V2
	c := &Client{
		UUID:    UUID(),
		handler: handler,
		version: v,
		tunnel:  transmitter.NewTunnel(wsConn, v),
		Online:  false,
	}
	//确认连接(握手,发送uuid)
	err = c.confirm(wsConn)
	if err != nil {
		return nil, err
	}
	//处理
	go c.handle()
	return c, nil
}

//
// confirm
// @Description:
// @receiver c
// @return error
//
func (c *Client) confirm(conn *transmitter.WSConn) error {
	return timer.TimeoutTask(func() error {
		_, err := conn.Write([]byte(c.UUID))
		if err != nil {
			return err
		}
		bytes := make([]byte, 32)
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		if string(bytes[:n]) != c.UUID {
			return errors.New("connection rejected")
		}
		c.Online = true
		return nil
	}, time.Second*10)
}

//
// handle
// @Description:
// @receiver c
//
func (c *Client) handle() {
	for {
		pl, err := c.tunnel.Read()
		if err == transmitter.ErrBadPacket {
			continue
		}
		if err != nil {
			//TODO 设置离线
			return
		}
		p := NewTransportPacket()
		err = p.Decode(pl)
		if err != nil {
			break
		}
		if p.UUID == c.UUID && p.Type != PacketTypeUnknown {
			switch p.Type {
			case PacketTypeLogin:
				//TODO 登录
			case PacketTypeLogout:
				//TODO 离线
			case PacketTypeMsg:
				c.handler.OnMessage(string(p.Payload))
			case PacketTypeReport:
				c.handler.OnReport(p.Payload)
			case PacketTypeKick:
				reply := AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, &reply)
				_ = log.Warn("server : ", reply.Message)
				//TODO 服务端控制离线
			case PacketTypeRestart:
				reply := AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, &reply)
				log.Info("server : ", reply.Message)
				//TODO 服务端控制重启
			}
		}
	}
}

//
// Disconnect
// @Description:
// @receiver c
// @param err
//
func (c *Client) Disconnect(err error) {
	//断开认证通道连接
	_ = c.tunnel.Close()
	//断开通信通道连接
	//事件必须在上层执行
	c.handler.OnDisconnect()
}
