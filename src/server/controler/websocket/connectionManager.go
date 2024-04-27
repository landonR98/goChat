package websocket

import (
	"encoding/json"
	"landonRyan/goChat/model"
	"landonRyan/goChat/util"
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	ActionSendMsg = "SEND_MSG"
	ActionNewMsg  = "NEW_MSG"
)

type wsMessage struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
}

type sendMsgPayload struct {
	ChatRoomId int    `json:"chatRoomId"`
	MessageTxt string `json:"messageTxt"`
}

type newMsgPayload struct {
	ChatRoomId int
	MessageTxt string
	Sender     model.User
}

var clients = make(map[int]*ws.Conn)

func WSHandler(res http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	user, _ := util.GetUserSession(req)
	clients[user.Id] = conn
	client := Client{UserId: user.Id, ChatRoomId: -1}
	store.addClient(client)
	for {
		mType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			continue
		}

		if mType == ws.TextMessage {
			var message wsMessage
			err := json.Unmarshal(msg, &message)
			if err != nil {
				log.Println(err)
				continue
			}

			switch message.Action {
			case ActionSendMsg:
				err = handleSendMessage(user.Id, message.Payload)
			}
		}
	}
}

func handleSendMessage(userId int, paylod string) error {
	var msgPayload sendMsgPayload
	err := json.Unmarshal([]byte(paylod), &msgPayload)
	if err != nil {
		return err
	}

	sendMsg := wsMessage{
		Action:  ActionNewMsg,
		Payload: paylod,
	}

	sendBytes, err := json.Marshal(sendMsg)

	client, err := store.getClientsByUserId(userId)
	if msgPayload.ChatRoomId != client.ChatRoomId {
		return err
	}
	sendClients, err := store.getClientsByChatRoomId(client.ChatRoomId)
	for _, sendClient := range sendClients {
		conn, found := clients[sendClient.UserId]
		if found {
			conn.WriteMessage(ws.TextMessage, sendBytes)
		} else {
			store.removeClientById(sendClient.UserId)
		}
	}

	return nil
}
