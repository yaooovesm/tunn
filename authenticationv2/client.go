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
	Online    bool                //是否在线
	version   transmitter.Version //传输版本
	tunnel    *transmitter.Tunnel //客户端连接
	sig       chan error          //同步信号
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
		//此时还没有建立传输连接，只需要关闭此处连接即可
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
			c.Disconnect(err)
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
				reply := &AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, reply)
				if reply.Ok {
					//TODO 登录成功
				} else {
					//登录失败
					c.Disconnect(errors.New(reply.Error))
				}
			case PacketTypeLogout:
				//无论如何都清除所有连接
				c.handler.OnLogout(nil)
			case PacketTypeMsg:
				c.handler.OnMessage(string(p.Payload))
			case PacketTypeReport:
				c.handler.OnReport(p.Payload)
			case PacketTypeKick:
				reply := AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, &reply)
				_ = log.Warn("server : ", reply.Message)
				c.handler.OnKick()
			case PacketTypeRestart:
				reply := AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, &reply)
				log.Info("server : ", reply.Message)
				c.handler.OnRestart()
			}
		}
	}
}

//
// Login
// @Description:
// @receiver c
// @return err
//
func (c *Client) Login() (err error) {
	configBytes, err := json.Marshal(config.Current)
	if err != nil {
		return ErrAuthBadPacket
	}
	p := TransportPacket{
		Type:    PacketTypeLogin,
		UUID:    c.UUID,
		Payload: configBytes,
	}
	_, err = c.tunnel.Write(p.Encode())
	if err != nil {
		return err
	}
	if err := timer.TimeoutTask(func() error {
		return <-c.sig
	}, time.Second*30); err != nil {
		log.Info("failed to login : ", err.Error())
		c.handler.OnLogin(err, nil)
		if err == timer.Timeout {
			return ErrAuthTimeout
		}
		return ErrAuthFailed
	}
	log.Info("login success")
	c.handler.OnLogin(nil, c.PublicKey)
	return nil
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
	c.handler.OnDisconnect(err)
	//设置离线
	c.Online = false
}
