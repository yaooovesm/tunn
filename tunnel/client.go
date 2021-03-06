package tunnel

import (
	"context"
	"errors"
	log "github.com/cihub/seelog"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"tunn/authenticationv2"
	"tunn/config"
	"tunn/config/protocol"
	"tunn/device"
	"tunn/networking"
	"tunn/traffic"
	"tunn/transmitter"
	"tunn/utils/timer"
)

//
// Client
// @Description:
//
type Client struct {
	IFace            device.Device
	Config           config.Config
	Context          context.Context
	Cancel           context.CancelFunc
	Error            error
	AuthClient       *authenticationv2.Client
	Address          string
	tunnelIndex      int
	maxIndex         int
	multiConn        *transmitter.MultiConn
	TxFP             *traffic.FlowProcessors
	Txfs             *traffic.FlowStatisticsFP
	RxFP             *traffic.FlowProcessors
	Rxfs             *traffic.FlowStatisticsFP
	mtu              int
	PK               []byte
	version          transmitter.Version
	handler          ClientConnHandler
	running          bool
	Online           bool
	txHandlerRunning bool
	SysRouteTable    *networking.SystemRouteTable
}

//
// NewClient
// @Description: 持久化多连接接模型
// @return *TCPClientV3
//
func NewClient() *Client {
	return &Client{
		IFace:            nil,
		Config:           config.Current,
		Error:            nil,
		tunnelIndex:      0,
		multiConn:        nil,
		mtu:              config.Current.Global.MTU,
		version:          transmitter.V2,
		running:          false,
		txHandlerRunning: false,
		Online:           false,
	}
}

//
// Init
// @Description:
// @receiver C
// @return error
//
func (c *Client) Init() error {
	//rx flow processor
	c.RxFP = traffic.NewFlowProcessor()
	c.RxFP.Name = "client_rx"
	//"RX : rx_packet_speed=", TXFs.PacketSpeed, "p/s rx_flow_speed=", TXFs.FlowSpeed/1024/1024, "mb/s"
	RXFs := &traffic.FlowStatisticsFP{Name: "rx"}
	c.RxFP.Register(RXFs, "rx_fs")
	c.Rxfs = RXFs
	//tx flow processor
	c.TxFP = traffic.NewFlowProcessor()
	c.TxFP.Name = "client_tx"
	//"TX : tx_packet_speed=", TXFs.PacketSpeed, "p/s tx_flow_speed=", TXFs.FlowSpeed/1024/1024, "mb/s"
	TXFs := &traffic.FlowStatisticsFP{Name: "tx"}
	c.TxFP.Register(TXFs, "tx_fs")
	c.Txfs = TXFs
	return nil
}

//
// Start
// @Description:
// @receiver C
// @return error
//
func (c *Client) Start(wg *sync.WaitGroup) error {
	wg.Add(1)
	//update key
	config.GenerateCipherKey()
	//multi conn
	c.multiConn = transmitter.NewMultiConn("client")
	//context
	ctx, cancelFunc := context.WithCancel(context.Background())
	c.Context = ctx
	c.Cancel = cancelFunc
	//auth
	//clientV3, err := authentication.NewClientV3(&AuthClientHandler{Client: c})
	client, err := authenticationv2.NewClient(&AuthClientHandler{Client: c})
	if err != nil {
		wg.Done()
		return err
	}
	c.AuthClient = client
	//connect
	err = c.AuthClient.Connect()
	if err != nil {
		wg.Done()
		return err
	}
	//login
	err = c.AuthClient.Login()
	if err != nil {
		wg.Done()
		return err
	}
	//验证服务器登录后覆盖本地配置需要重新载入
	//重新载入配置
	c.Config = config.Current
	c.mtu = config.Current.Global.MTU
	//preprocess Address
	c.Address = strings.Join([]string{c.Config.Global.Address, strconv.Itoa(c.Config.Global.Port)}, ":")

	//设置传输协议
	switch c.Config.Global.Protocol {
	case protocol.TCP:
		c.handler = &TCPClientHandler{}
	case protocol.KCP:
		c.handler = &KCPClientHandler{}
	case protocol.WS:
		c.handler = &WSClientHandler{}
	case protocol.WSS:
		c.handler = &WSSClientHandler{}
	default:
		wg.Done()
		return errors.New("unsupported protocol : " + string(c.Config.Global.Protocol))
	}
	log.Info("transmit protocol : ", c.Config.Global.Protocol.ToString())

	//get interface cidr address after login
	//iface
	if c.IFace == nil {
		dev, err := device.NewTunDevice()
		if err != nil {
			wg.Done()
			return err
		}
		err = dev.Setup()
		if err != nil {
			wg.Done()
			return err
		}
		//注册系统路由表
		c.SysRouteTable = networking.NewSystemRouteTable(dev.Name())
		c.IFace = dev
	} else {
		//客户端可能重置网卡IP
		err := c.IFace.OverwriteCIDR(config.Current.Device.CIDR)
		if err != nil {
			wg.Done()
			return err
		}
	}
	//更新系统路由表
	c.SysRouteTable.Merge(config.Current.Routes)
	c.SysRouteTable.DeployAll()
	//同步size
	size := config.Current.Global.MultiConn
	log.Info("multi connection size : ", size)
	//初始化完成
	c.handler.AfterInitialize(c)
	for i := 0; i < size; i++ {
		conn, err := c.handler.CreateAndSetup(c.Address, c.Config)
		if err != nil {
			wg.Done()
			return err
		}
		//confirm
		err = c.confirm(conn)
		if err != nil {
			wg.Done()
			return err
		}
		num := c.multiConn.Push(conn)
		go c.RXHandler(conn, num)
	}
	c.running = true
	c.Online = true
	if !c.txHandlerRunning {
		go c.TXHandler()
	}
	log.Info("connected to the server successfully!")
	log.Info("your ip address is ", config.Current.Device.CIDR, ".")
	wg.Done()
	select {
	case <-c.Context.Done():
		c.running = false
		c.Online = false
		err := c.Error
		c.Error = nil
		return err
	}
}

//
// confirm
// @Description: 验证通道UUID
// @receiver c
// @param conn
//
func (c *Client) confirm(conn net.Conn) error {
	return timer.TimeoutTask(func() error {
		uuid := c.AuthClient.UUID
		log.Info("send confirm packet : ", uuid)
		_, err := conn.Write([]byte(uuid))
		if err != nil {
			return err
		}
		bytes := make([]byte, 32)
		n, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		if uuid != string(bytes[:n]) {
			return errors.New("connection rejected")
		}
		return nil
	}, time.Second*10)
}

//
// Terminate
// @Description:
// @receiver c
//
func (c *Client) Terminate() {
	c.running = false
	if c.AuthClient != nil {
		_ = c.AuthClient.Logout()
	}
	c.Error = errors.New("terminated")
	c.Cancel()
}

//
// Stop
// @Description:
// @receiver C
//
func (c *Client) Stop() {
	c.running = false
	c.Cancel()
}

//
// SetErr
// @Description:
// @receiver c
// @param err
//
func (c *Client) SetErr(err error) {
	c.Error = err
}

//
// TXHandler
// @Description:
// @receiver c
//
func (c *Client) TXHandler() {
	defer func() {
		log.Info("[tx_handler] exit")
		c.txHandlerRunning = false
	}()
	c.txHandlerRunning = true
	packet := make([]byte, c.Config.Global.MTU)
	for {
		//interface -> tunnel
		n, err := c.IFace.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		//流量处理
		if c.running {
			if conn := c.multiConn.Get(); conn != nil {
				//先处理流量再封包
				_, _ = conn.Write(c.TxFP.Process(packet[:n]))
			}
		}
	}
}

//
// RXHandler
// @Description:
// @receiver C
//
func (c *Client) RXHandler(conn net.Conn, num int) {
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	//封包器
	reader := transmitter.NewTunReader(conn, c.version)
	for {
		pl, err := reader.Read()
		if err != nil && err != transmitter.ErrBadPacket {
			log.Info("[rx][#", num, "] exit : ", err.Error())
			if c.running && c.Error == nil {
				c.SetErr(ErrDisconnectAccidentally)
				c.Stop()
			}
			return
		}
		//流量处理
		_, _ = c.IFace.Write(c.RxFP.Process(pl))
	}
}
