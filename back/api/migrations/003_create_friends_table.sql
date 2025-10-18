-- friendsテーブルの作成

CREATE TABLE IF NOT EXISTS friends (
    id CHAR(36) PRIMARY KEY COMMENT 'フレンドリクエストID (UUID)',
    requester_id CHAR(36) NOT NULL COMMENT '申請者のユーザーID',
    addressee_id CHAR(36) NOT NULL COMMENT '申請相手のユーザーID',
    status ENUM('pending', 'accepted', 'rejected') NOT NULL DEFAULT 'pending' COMMENT 'ステータス',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    -- インデックス
    INDEX idx_requester_id (requester_id),
    INDEX idx_addressee_id (addressee_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),

    -- ユニーク制約（同じペアのリクエストは1つのみ）
    UNIQUE KEY unique_friendship (requester_id, addressee_id),

    -- 外部キー制約
    CONSTRAINT fk_friends_requester_id
        FOREIGN KEY (requester_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_friends_addressee_id
        FOREIGN KEY (addressee_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='フレンドテーブル';
