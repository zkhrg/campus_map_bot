package domain

import "time"

type Peer struct {
	Row     string    `json:"row"`
	Col     string    `json:"col"`
	Cluster string    `json:"cluster"`
	Time    time.Time `json:"time"`
	Status  string    `json:"status"`
}

type FriendsByTelegramID struct {
	TgID     int64  `json:"tg_id"`
	PeerName string `json:"peer_name"`
}

type PeerStatusResponse struct {
	PeerName string    `json:"peer_name"`
	Row      string    `json:"row"`
	Col      string    `json:"col"`
	Cluster  string    `json:"cluster"`
	Time     time.Time `json:"time"`
	Status   string    `json:"status"`
}

type UpdatePeersRequest struct {
	Peers []struct {
		PeerName string `json:"peer_name"`
		Row      string `json:"row"`
		Col      string `json:"col"`
		Cluster  string `json:"cluster"`
		Status   string `json:"status"`
	} `json:"peers"`
}

type AddFriendRequest struct {
	TgID     int64  `json:"tg_id"`
	PeerName string `json:"peer_name"`
}

type DeleteFriendRequest struct {
	TgID     int64  `json:"tg_id"`
	PeerName string `json:"peer_name"`
}
