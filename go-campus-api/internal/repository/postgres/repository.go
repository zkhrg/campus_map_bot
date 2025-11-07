package postgres

import (
	"context"
	"time"

	"go-campus-api/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
	psql sq.StatementBuilderType
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) GetPeerStatus(ctx context.Context, peerName string) (*domain.PeerStatusResponse, error) {
	query, args, err := r.psql.
		Select("peer_name", "row", "col", "cluster", "time", "status").
		From("peers").
		Where(sq.Eq{"peer_name": peerName}).
		ToSql()

	if err != nil {
		return nil, err
	}

	var peer domain.PeerStatusResponse
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&peer.PeerName,
		&peer.Row,
		&peer.Col,
		&peer.Cluster,
		&peer.Time,
		&peer.Status,
	)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	return &peer, err
}

func (r *Repository) UpdatePeers(ctx context.Context, peers []domain.Peer, peerNames []string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for i, peer := range peers {
		query, args, err := r.psql.
			Insert("peers").
			Columns("peer_name", "row", "col", "cluster", "time", "status").
			Values(peerNames[i], peer.Row, peer.Col, peer.Cluster, time.Now(), "1").
			Suffix("ON CONFLICT (peer_name) DO UPDATE SET row = EXCLUDED.row, col = EXCLUDED.col, cluster = EXCLUDED.cluster, time = EXCLUDED.time, status = EXCLUDED.status").
			ToSql()

		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *Repository) GetFriendsStatus(ctx context.Context, tgID int64) ([]domain.PeerStatusResponse, error) {
	query, args, err := r.psql.
		Select("p.peer_name", "p.row", "p.col", "p.cluster", "p.time", "p.status").
		From("friends f").
		Join("peers p ON f.peer_name = p.peer_name").
		Where(sq.Eq{"f.tg_id": tgID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []domain.PeerStatusResponse
	for rows.Next() {
		var peer domain.PeerStatusResponse
		err := rows.Scan(&peer.PeerName, &peer.Row, &peer.Col, &peer.Cluster, &peer.Time, &peer.Status)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}

	return peers, rows.Err()
}

func (r *Repository) AddFriend(ctx context.Context, tgID int64, peerName string) error {
	query, args, err := r.psql.
		Insert("friends").
		Columns("tg_id", "peer_name").
		Values(tgID, peerName).
		Suffix("ON CONFLICT (tg_id, peer_name) DO NOTHING").
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}

func (r *Repository) DeleteFriend(ctx context.Context, tgID int64, peerName string) error {
	query, args, err := r.psql.
		Delete("friends").
		Where(sq.Eq{"tg_id": tgID, "peer_name": peerName}).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query, args...)
	return err
}
