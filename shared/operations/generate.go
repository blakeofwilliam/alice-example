package operations

import (
	"fmt"
	"log"

	"github.com/blakeofwilliam/alice-example/shared/peerManager"
	"github.com/getamis/alice/crypto/ecpointgrouplaw"
	"github.com/getamis/alice/crypto/elliptic"
	"github.com/getamis/alice/crypto/tss/ecdsa/cggmp/dkg"
	"github.com/getamis/alice/types"
	"google.golang.org/protobuf/proto"
)

type PublicKey struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type ECPoint struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type PaillierKey struct {
	P string `json:"p"`
	Q string `json:"q"`
}

type PederssenOpenParameter struct {
	N string `json:"n"`
	S string `json:"s"`
	T string `json:"t"`
}

type BK struct {
	X    string `json:"x"`
	Rank uint32 `json:"rank"`
}

type GenerateResult struct {
	BKs               map[string]BK     				`json:"bks"`
	PartialPublicKeys map[string]PublicKey 			`json:"partialPublicKeys"`
	PublicKey         PublicKey            			`json:"publicKey"`
	Share             string                   	`json:"share"`
}

type GenerateService struct {
	types.StateChangedListener

	dkg 					*dkg.DKG
	errorChannel 	chan error
	peerManager 	*peerManager.PeerManager
}

var Curves = map[string]elliptic.Curve{
	"SECP256K1": elliptic.Secp256k1(),
	"ED25519": elliptic.Ed25519(),
}

func GetCurve() elliptic.Curve {
	return elliptic.Secp256k1()
}

func NewGenerateService(curveId string, peerManager *peerManager.PeerManager) (service *GenerateService, err error) {
	curve := Curves[curveId]

	if peerManager.SelfID() == "server" {
		peerManager.AddPeer("client")
	} else {
		peerManager.AddPeer("server")
	}

	service = &GenerateService{
		errorChannel: make(chan error, 1),
		peerManager: peerManager,
	}

	sid := []byte("server-client")
	dkg, err := dkg.NewDKG(curve, peerManager, sid, 2, 0, service)
	if err == nil {
		service.dkg = dkg
	} else {
		log.Println("<<< Failed to create DKG", err)
		return
	}

	return
}

func (service *GenerateService) Done() <-chan error {
	return service.errorChannel
}

func (service *GenerateService) Listen() {
	for {
		data := &dkg.Message{}
		websocket := service.peerManager.WebSocket()

		_, message, err := websocket.Socket.ReadMessage()
		if err != nil {
			log.Println("<<< Error reading from socket: ", string(err.Error()))
			service.errorChannel <- err
			return
		}

		err = proto.Unmarshal(message, data)
		if err != nil {
			log.Println("<<< Error unmarshaling message", err)
			service.errorChannel <- err
			return
		}

		err = service.dkg.AddMessage(data.GetId(), data)
		if err != nil {
			log.Println("<<< Error adding message to DKG: ", err)
			service.errorChannel <- err
			return
		}

		log.Println("<<< Added message to DKG: ", data.String())

		state := service.dkg.GetState()
		log.Println("<<< DKG State: ", state.String())

		if data.Type == dkg.Type_Result {
			log.Println("<<< Success")
			service.errorChannel <- nil
			return
		}
	}
}

func (service *GenerateService) OnStateChanged(oldState types.MainState, newState types.MainState) {
	log.Println("State changed", "old", oldState.String(), "new", newState.String())

	if newState == types.StateFailed {
		service.errorChannel <- fmt.Errorf("state %s -> %s", oldState.String(), newState.String())
		return
	} else if newState == types.StateDone {
		service.errorChannel <- nil
		return
	}
}

func (service *GenerateService) Process() (result *GenerateResult, err error) {
	go service.Listen()

	service.dkg.Start()
	defer service.dkg.Stop()

	select {
		case err = <-service.errorChannel:
			if err != nil {
				return
			}
		default:
			// No errors, all is well
	}

	if err = <-service.Done(); err != nil {
		return
	}

	dkgResult, err := service.dkg.GetResult()
	if err != nil {
		return
	}

	log.Println("!!! Got DKG result: ", dkgResult)

	result = service.convertDkgResult(dkgResult)

	return
}

func (service *GenerateService) convertDkgResult(dkgResult *dkg.Result) (result *GenerateResult) {
	log.Println("<<< Converting DKG result: ", dkgResult)
	bks := make(map[string]BK)
	bk := &BK{
		X: dkgResult.Bks[service.peerManager.SelfID()].GetX().String(), 
		Rank: dkgResult.Bks[service.peerManager.SelfID()].GetRank(),
	}
	bks[service.peerManager.SelfID()] = *bk
	for _, peer := range service.peerManager.PeerIDs() {
		bk = &BK{
			X: dkgResult.Bks[peer].GetX().String(), 
			Rank: dkgResult.Bks[peer].GetRank(),
		}
		bks[peer] = *bk
	}

	var partialPublicKeys = map[string]*ecpointgrouplaw.ECPoint{}
	partialPublicKey := ecpointgrouplaw.ScalarBaseMult(GetCurve(), dkgResult.Share)
	partialPublicKeys[service.peerManager.SelfID()] = partialPublicKey

	result = &GenerateResult{
		bks, 
		make(map[string]PublicKey),
		PublicKey{
			X: dkgResult.PublicKey.ToPubKey().X.String(), 
			Y: dkgResult.PublicKey.GetY().String(),
		}, 
		dkgResult.Share.String(),
	}

	for peerId, partialPublicKey := range partialPublicKeys {
		result.PartialPublicKeys[peerId] = PublicKey{
			X: partialPublicKey.GetX().String(),
			Y: partialPublicKey.GetY().String(),
		}
	}

	return 
}
