package core

import (
	"encoding/json"
	"html"
	"log"
	"strconv"
	"time"
)

type Message struct {

	//"type":"say|siliao","to_client_id":"all","to_client_name":"所有人","content"
	Type         string `json:"type"`
	ToClinetId   string `json:"to_client_id"`
	ToClientName string `json:"to_client_name"`
	Content      string `json:"content"`
	ClientName   string `json:"client_name"`
	RoomId       int    `json:"room_id"`
	ClinetId     int    `json:"client_id"`
}

/**
*消息处理结构体
 */
type MessageHandle struct {
	Cli  *Client
	Data []byte
}

func NewMessageHandle(cli *Client, data []byte) *MessageHandle {

	return &MessageHandle{Cli: cli, Data: data}
}

/**
*消息解析及逻辑处理
 */
func (this *MessageHandle) Handled() {
	defer func() {
		recover()
	}()
	var msg Message
	if err := json.Unmarshal(this.Data, &msg); err != nil {
		log.Printf("\n json.Unmarshal err:%v \n", err)
		return
	}
	hData := make(map[string]interface{})
	switch msg.Type {
	case "login":
		this.Cli.ClientName = msg.ClientName
		hData["type"] = "login"
		hData["client_name"] = html.EscapeString(msg.ClientName)
		hData["client_id"] = this.Cli.ClientId
		hData["client_list"] = this.GetClientList()
		hData["time"] = time.Now().Format(DateFormat)

	case "say":
		//data['from_client_id'], data['from_client_name'], data['content'], data['time']
		hData["type"] = "say"
		hData["from_client_id"] = this.Cli.ClientId
		hData["from_client_name"] = this.Cli.ClientName
		hData["to_client_id"] = msg.ToClinetId
		hData["to_client_name"] = msg.ToClientName
		hData["content"] = html.EscapeString(msg.Content)
		hData["time"] = time.Now().Format(DateFormat)
	case "siliao":
		var toClient *Client
		if toid, err2 := strconv.Atoi(msg.ToClinetId); err2 == nil {
			toClient = this.GetClientById(toid)
			if toClient == nil {
				log.Println("接收对象不存在或已退出")
				return
			}
		} else {
			log.Println("ToClinetId错误")
			return
		}
		hData["type"] = "siliao"
		hData["from_client_id"] = this.Cli.ClientId
		hData["from_client_name"] = this.Cli.ClientName
		hData["to_client_id"] = msg.ToClinetId
		hData["to_client_name"] = msg.ToClientName
		hData["content"] = html.EscapeString(msg.Content)
		hData["time"] = time.Now().Format(DateFormat)
		if sendData, err2 := json.Marshal(hData); err2 != nil {
			log.Printf("\n json.Marshal err:%v \n", err2)
			return
		} else {
			toClient.send <- sendData
			this.Cli.send <- sendData
		}
		return
	}
	if sendData, err2 := json.Marshal(hData); err2 != nil {
		log.Printf("\n json.Marshal err:%v \n", err2)
		return
	} else {
		this.Cli.hub.broadcast <- sendData
	}
}

/**
*获取所有用户列表
 */
func (this *MessageHandle) GetClientList() map[int]string {
	m := make(map[int]string)
	for k, _ := range this.Cli.hub.clients {
		m[k.ClientId] = k.ClientName
	}
	return m
}

/**
*根据clinentId获取client
 */
func (this *MessageHandle) GetClientById(id int) *Client {
	for k, _ := range this.Cli.hub.clients {
		if k.ClientId == id {
			return k
		}
	}
	return nil
}
