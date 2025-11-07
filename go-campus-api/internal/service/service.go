package service

import (
	"context"
	"go-campus-api/internal/domain"
)

type Repository interface {
	GetPeerStatus(ctx context.Context, peerName string) (*domain.PeerStatusResponse, error)
	UpdatePeers(ctx context.Context, peers []domain.Peer, peerNames []string) error
	GetFriendsStatus(ctx context.Context, tgID int64) ([]domain.PeerStatusResponse, error)
	AddFriend(ctx context.Context, tgID int64, peerName string) error
	DeleteFriend(ctx context.Context, tgID int64, peerName string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetPeerStatus(ctx context.Context, peerName string) (*domain.PeerStatusResponse, error) {
	return s.repo.GetPeerStatus(ctx, peerName)
}

func (s *Service) UpdatePeers(ctx context.Context, req domain.UpdatePeersRequest) error {
	peers := make([]domain.Peer, len(req.Peers))
	peerNames := make([]string, len(req.Peers))

	for i, p := range req.Peers {
		peers[i] = domain.Peer{
			Row:     p.Row,
			Col:     p.Col,
			Cluster: p.Cluster,
			Status:  p.Status,
		}
		peerNames[i] = p.PeerName
	}

	return s.repo.UpdatePeers(ctx, peers, peerNames)
}

func (s *Service) GetFriendsStatus(ctx context.Context, tgID int64) ([]domain.PeerStatusResponse, error) {
	return s.repo.GetFriendsStatus(ctx, tgID)
}

func (s *Service) AddFriend(ctx context.Context, tgID int64, peerName string) error {
	return s.repo.AddFriend(ctx, tgID, peerName)
}

func (s *Service) DeleteFriend(ctx context.Context, tgID int64, peerName string) error {
	return s.repo.DeleteFriend(ctx, tgID, peerName)
}
