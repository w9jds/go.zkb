package zkb

import (
	"encoding/json"
	"fmt"
	"log"

	esi "github.com/w9jds/go.esi"
)

// ZkillResponse struct to be returned when calling the zkill api
type ZkillResponse struct {
	ID  uint32     `json:"killmail_id,omitempty"`
	Zkb zkillboard `json:"zkb,omitempty"`
}

// RedisResponse struct to be returned when calling the redis zkill queue
type RedisResponse struct {
	ID       uint32       `json:"killID,omitempty"`
	Zkb      zkillboard   `json:"zkb,omitempty"`
	Killmail esi.KillMail `json:"killmail,omitempty"`
}

type redisPackage struct {
	Package RedisResponse `json:"package"`
}

type zkillboard struct {
	LocationID     uint32   `json:"locationID,omitempty"`
	Hash           string   `json:"hash,omitempty"`
	FittedValue    float64  `json:"fittedValue,omitempty"`
	DroppedValue   float64  `json:"droppedValue,omitempty"`
	DestroyedValue float64  `json:"destroyedValue,omitempty"`
	TotalValue     float64  `json:"totalValue,omitempty"`
	NPC            bool     `json:"npc,omitempty"`
	Solo           bool     `json:"solo,omitempty"`
	Awox           bool     `json:"awox,omitempty"`
	Labels         []string `json:"labels,omitempty"`
	Href           string   `json:"href,omitempty"`
}

var zkillBase = "https://zkillboard.com/api"
var redisqBase = "https://redisq.zkillboard.com"

// GetKillMail returns the basic information about the killmail from zkill
func (zkb *Client) GetKillMail(killID string) (*ZkillResponse, error) {
	body, error := zkb.get(zkillBase, fmt.Sprintf("/killID/%s/", killID))
	if error != nil {
		return nil, error
	}

	var killmail []ZkillResponse
	error = json.Unmarshal(body, &killmail)
	if error != nil {
		return nil, error
	}

	return &killmail[0], nil
}

// GetRedisItem returns the next item in the redis queue from zKill
func (zkb *Client) GetRedisItem(queueID string) (*RedisResponse, error) {
	body, error := zkb.get(redisqBase, fmt.Sprintf("/listen.php?ttw=0&queueID=%s", queueID))
	if error != nil {
		log.Println("Error creating new request: ", error)
		return nil, error
	}

	var bundle redisPackage
	error = json.Unmarshal(body, &bundle)
	if error != nil {
		log.Println("Error parsing bundle response")
		return nil, error
	}

	return &bundle.Package, nil
}
