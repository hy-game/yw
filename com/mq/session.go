package mq

import (
	"com/log"
	"com/util"
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
	"math"
	"pb"
	"time"
)

//NsqSession	服务器间消息会话
type NsqSession struct {
	handler  *routeNsq
	producer *nsq.Producer
	cfg      *nsq.Config
	lkupAddr []string
}

func NewNsqSession(nsqdAddr string, nsqLkup []string) *NsqSession {
	m := &NsqSession{
		handler: newRouteCt(math.MaxUint16),
	}

	m.cfg = nsq.NewConfig()
	m.cfg.LookupdPollInterval = time.Second * 5 //设置重连时间

	//producer
	m.addProducer(nsqdAddr)

	//consumer
	m.lkupAddr = nsqLkup

	return m
}

//Send 发送消息
func (this *NsqSession) Send(topic string, msgId pb.MsgIDS2S, msgData proto.Message, serId uint16) error {
	data, err := WritePkg(uint16(msgId), msgData, serId)
	if err != nil {
		return err
	}
	err = this.producer.Publish(topic, data)
	if err != nil {
		return err
	}
	if msgId > pb.MsgIDS2S_S2STraceEnd {
		var msgStr string
		if msgData != nil {
			msgStr = msgData.String()
		}
		log.Tracef("send to %s [%d][%s] %s", topic, msgId, pb.MsgIDS2S_name[int32(msgId)], msgStr)
	}
	return nil
}

//Stop	停止
func (this *NsqSession) Stop() {
	if this.producer != nil {
		this.producer.Stop()
	}
}

//RegisterMsgHandle	注册消息处理函数
func (this *NsqSession) RegisterMsgHandle(msgID pb.MsgIDS2S, cf func() proto.Message, df func(msg proto.Message, serId uint16)) {
	this.handler.register(msgID, cf, df)
}

//HandleMessage 消息处理入口，nsq自己调用
func (this *NsqSession) HandleMessage(msg *nsq.Message) error {
	defer func() {
		if err := recover(); err != nil {
			log.Error(err)
			util.PrintStack()
		}
	}()

	msgId, data, serId := ParserPkg(msg.Body)
	err := this.handler.handle(msgId, data, serId)
	if err != nil && msgId > uint16(pb.MsgIDS2S_S2STraceEnd) {
		log.Warnf("handle msg %d err:%v", msgId, err)
	}
	return nil
}

//AddConsumer	添加一个消费者
func (this *NsqSession) AddConsumer(topic string, chl string) {
	consumer, err := nsq.NewConsumer(topic, chl, this.cfg)
	if err != nil {
		log.Panic(err)
		return
	}

	consumer.AddHandler(this)
	consumer.SetLoggerLevel(nsq.LogLevelWarning)

	err = consumer.ConnectToNSQLookupds(this.lkupAddr)
	if err != nil {
		log.Panic(err)
	}
	log.Infof("add consumer t:%s, c:%s to %s", topic, chl, this.lkupAddr)
}

//addProducer	添加一个生产者
func (this *NsqSession) addProducer(addr string) {
	var err error
	this.producer, err = nsq.NewProducer(addr, this.cfg)
	if err != nil {
		log.Panic(err)
	}
	this.producer.SetLoggerLevel(nsq.LogLevelWarning)
	log.Infof("add producer to %s", addr)
}
