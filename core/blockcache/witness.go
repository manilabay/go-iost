package blockcache

import (
	"encoding/json"
	"errors"

	"github.com/iost-official/go-iost/db"
	"github.com/iost-official/go-iost/vm/database"
)

// SetPending set pending witness list
func (wl *WitnessList) SetPending(pl []string) {
	wl.PendingWitnessList = pl
}

// SetActive set active witness list
func (wl *WitnessList) SetActive(al []string) {
	wl.ActiveWitnessList = al
}

// Pending get pending witness list
func (wl *WitnessList) Pending() []string {
	return wl.PendingWitnessList
}

// Active get active witness list
func (wl *WitnessList) Active() []string {
	return wl.ActiveWitnessList
}

// NetID get net id
func (wl *WitnessList) NetID() []string {
	return wl.WitnessInfo
}

// UpdatePending update pending witness list
func (wl *WitnessList) UpdatePending(mv db.MVCCDB) error {

	vi := database.NewVisitor(0, mv)

	jwl := database.MustUnmarshal(vi.Get("vote_producer.iost-" + "pendingProducerList"))
	if jwl == nil {
		return errors.New("failed to get pending list")
	}
	str := make([]string, 0)
	err := json.Unmarshal([]byte(jwl.(string)), &str)
	if err != nil {
		return err
	}
	wl.SetPending(str)

	return nil
}

// UpdateInfo update pending witness list
func (wl *WitnessList) UpdateInfo(mv db.MVCCDB) error {

	wl.WitnessInfo = make([]string, 0, 0)
	vi := database.NewVisitor(0, mv)
	for _, v := range wl.PendingWitnessList {
		iAcc := database.MustUnmarshal(vi.MGet("vote_producer.iost-producerKeyToId", v))
		if iAcc == nil {
			continue
		}

		var acc []string
		err := json.Unmarshal([]byte(iAcc.(string)), &acc)
		if err != nil {
			continue
		}

		jwl := database.MustUnmarshal(vi.MGet("vote_producer.iost-producerTable", acc[0]))
		if jwl == nil {
			continue
		}

		var str WitnessInfo
		err = json.Unmarshal([]byte(jwl.(string)), &str)
		if err != nil {
			continue
		}
		wl.WitnessInfo = append(wl.WitnessInfo, str.NetID)
	}
	return nil
}

// CopyWitness is copy witness
func (wl *WitnessList) CopyWitness(n *BlockCacheNode) {
	if n == nil {
		return
	}
	wl.SetActive(n.Active())
	wl.SetPending(n.Pending())
}
