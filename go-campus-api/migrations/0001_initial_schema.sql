-- +goose Up
CREATE TABLE IF NOT EXISTS peers (
    peer_name VARCHAR(16) PRIMARY KEY,
    row VARCHAR(50) NOT NULL DEFAULT '',
    col VARCHAR(50) NOT NULL DEFAULT '',
    cluster VARCHAR(100) NOT NULL DEFAULT '',
    time TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(10) NOT NULL DEFAULT '0'
);

CREATE TABLE IF NOT EXISTS friends (
    tg_id BIGINT NOT NULL,
    peer_name VARCHAR(16) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tg_id, peer_name),
    FOREIGN KEY (peer_name) REFERENCES peers(peer_name) ON DELETE CASCADE
);

CREATE INDEX idx_friends_tg_id ON friends(tg_id);
CREATE INDEX idx_peers_status ON peers(status);

-- +goose Down
DROP TABLE IF EXISTS friends;
DROP TABLE IF EXISTS peers;