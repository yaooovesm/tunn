package authenticationv2

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
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
	"tunn/networking"
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
	//创建客户端
	c := &Client{
		UUID:    UUID(),
		handler: handler,
		version: transmitter.V2,
		Online:  false,
		sig:     make(chan error, 1),
	}
	return c, nil
}

//
// Connect
// @Description:
// @receiver c
// @return *transmitter.WSConn
// @return error
//
func (c *Client) Connect() error {
	//设置连接地址
	address := strings.Join([]string{config.Current.Auth.Address, strconv.Itoa(config.Current.Auth.Port)}, ":")
	log.Info("connect to authentication server : ", address)
	//读取证书
	pool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(config.Current.Security.CertPem)
	if err != nil {
		return log.Error("cannot load cert : ", err)
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
		return ErrConnectFailed
	}
	wsConn := transmitter.WrapWSConn(conn)
	//设置连接
	c.tunnel = transmitter.NewTunnel(wsConn, c.version)
	//确认连接(握手,发送uuid)
	err = c.confirm(wsConn)
	if err != nil {
		return err
	}
	//处理
	go c.handle()
	return nil
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
Loop:
	for {
		pl, err := c.tunnel.Read()
		if err == transmitter.ErrBadPacket {
			continue
		}
		if err != nil {
			c.Disconnect(err, true)
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
				reply := &AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, reply)
				if reply.Ok {
					err := c.afterLogin(reply)
					if err != nil {
						//处理失败
						c.sig <- err
						//退出循环
						break Loop
					}
				} else {
					//登录失败
					c.Disconnect(errors.New(reply.Error), true)
					//退出循环
					break Loop
				}
				c.sig <- nil
			case PacketTypeLogout:
				c.sig <- nil
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
				//断开当前客户端
				c.Disconnect(nil, false)
			case PacketTypeRestart:
				reply := AuthReply{}
				//解析reply
				_ = json.Unmarshal(p.Payload, &reply)
				log.Info("server : ", reply.Message)
				c.handler.OnRestart()
				//断开当前客户端,重启后将会重新创建客户端
				c.Disconnect(nil, false)
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
	//等待信号
	if err := timer.TimeoutTask(func() error {
		return <-c.sig
	}, time.Second*30); err != nil {
		log.Info("failed to login : ", err.Error())
		c.handler.OnLogin(err, nil)
		if err == timer.Timeout {
			return ErrAuthTimeout
		}
		return err
	}
	log.Info("login success")
	c.handler.OnLogin(nil, c.PublicKey)
	return nil
}

//
// afterLogin
// @Description:
// @receiver c
// @param reply
// @return error
//
func (c *Client) afterLogin(reply *AuthReply) error {
	data := make(map[string]string)
	//解码服务端发送的配置文件
	err := json.Unmarshal([]byte(reply.Message), &data)
	if err != nil {
		return err
	}
	//接收ws_key
	if wskey, ok := data["ws_key"]; ok && wskey != "" {
		c.WSKey = wskey
	}
	//接收配置
	if cfg, ok := data["config"]; ok && cfg != "" {
		pushedConfig := config.PushedConfig{}
		err := json.Unmarshal([]byte(cfg), &pushedConfig)
		if err != nil {
			return errors.New("failed to fetch config")
		}
		if pushedConfig.Device.CIDR == "" {
			return errors.New("failed to get a address")
		}
		//覆盖配置到本地
		config.Current.MergePushed(pushedConfig)
		for i := range pushedConfig.Routes {
			//当有网络暴露时
			if pushedConfig.Routes[i].Option == config.RouteOptionExport {
				//开启路由支持
				networking.RouteSupport()
				break
			}
		}
	} else {
		return errors.New("failed to fetch config")
	}
	//接收公钥
	if key, ok := data["key"]; ok && key != "" {
		keyBytes, err := hex.DecodeString(key)
		if err != nil {
			return err
		}
		c.PublicKey = keyBytes
		log.Info("receive ", len(c.PublicKey), " bytes key from server")
	}
	return nil
}

func (c *Client) Logout() error {
	if !c.Online {
		return errors.New("client offline")
	}
	//读取当前配置文件
	configBytes, err := json.Marshal(config.Current)
	if err != nil {
		return ErrAuthBadPacket
	}
	//发送包
	p := TransportPacket{
		Type:    PacketTypeLogout,
		UUID:    c.UUID,
		Payload: configBytes,
	}
	_, err = c.tunnel.Write(p.Encode())
	if err != nil {
		return ErrAuthBadPacket
	}
	//等待5秒服务器响应
	_ = timer.TimeoutTask(func() error {
		return <-c.sig
	}, time.Second*5)
	//无论响应如何都关闭连接
	//退出登录事件
	c.handler.OnLogout(nil)
	//清除登录
	c.Disconnect(nil, false)
	log.Info("logout success")
	return nil
}

//
// Disconnect
// @Description:
// @receiver c
// @param err
//
func (c *Client) Disconnect(err error, event bool) {
	//断开认证通道连接
	_ = c.tunnel.Close()
	//断开通信通道连接
	//事件必须在上层执行
	if event {
		c.handler.OnDisconnect(err)
	}
	//设置离线
	c.Online = false
}

//
// Report
// @Description:
// @receiver c
// @param data
// @return error
//
func (c *Client) Report(data []byte) (err error) {
	p := TransportPacket{
		Type:    PacketTypeReport,
		UUID:    c.UUID,
		Payload: data,
	}
	_, err = c.tunnel.Write(p.Encode())
	return
}

//
// Message
// @Description:
// @receiver c
// @param msg
// @return error
//
func (c *Client) Message(msg string) (err error) {
	p := TransportPacket{
		Type:    PacketTypeMsg,
		UUID:    c.UUID,
		Payload: []byte(msg),
	}
	_, err = c.tunnel.Write(p.Encode())
	return
}
