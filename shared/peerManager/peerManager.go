package peerManager

import (
	"log"

	"github.com/blakeofwilliam/alice-example/shared/websocket"
	"github.com/getamis/alice/types"
	"google.golang.org/protobuf/proto"
)

type PeerManager struct {
	types.PeerManager
	
	id					string
	peers				[]string
	protocol 	  string
	websocket		*websocket.WebSocket
}

func NewPeerManager(
	id string,
	socket *websocket.WebSocket,
	protocol string,
) *PeerManager {
	return &PeerManager{
		id: id,
		peers: []string{},
		protocol: protocol,
		websocket: socket,
	}
}

func (peerManager *PeerManager) AddPeer(id string) {
	peerManager.peers = append(peerManager.peers, id)
}

func (peerManager *PeerManager) MustSend(peerID string, message interface{}) {
	log.Println("<<< MustSend message to peerID: ", peerID, " message: ", message)
	peerManager.send(message)
}

func (peerManager *PeerManager) NumPeers() uint32 {
	return uint32(len(peerManager.peers))
}

func (peerManager *PeerManager) PeerIDs() []string {
	return peerManager.peers
}

func (peerManager *PeerManager) SelfID() string {
	return peerManager.id
}

func (peerManager *PeerManager) WebSocket() *websocket.WebSocket {
	return peerManager.websocket
}

/*******************************************
 * Private Methods
 *******************************************/

func (peerManager *PeerManager) send(message interface{}) {
	// Handle errors
	defer func() {
		panicErr := recover()
		if panicErr != nil {
			log.Println("<<< Hit a panic error in send: ", panicErr)
		}
	}()

	// Handle lock management
	// peerManager.lock.Lock()
	// defer peerManager.lock.Unlock()

	msg, ok := message.(proto.Message)
	if !ok {
		log.Println("<<< Network Error - invalid proto message to send")
		return 
	}

	byteString, err := proto.Marshal(msg)
	if err != nil {
		log.Println("<<< Cannot marshal message: " + err.Error())
		return 
	}

	err = peerManager.websocket.Write(byteString)
	if err != nil {
		log.Println("<<< Network Error - " + err.Error())
		return 
	}
}