-- group_membersテーブルの作成

CREATE TABLE IF NOT EXISTS group_members (
    id CHAR(36) NOT NULL COMMENT 'メンバーID (UUID)',
    group_id CHAR(36) NOT NULL COMMENT 'グループID',
    user_id CHAR(36) NOT NULL COMMENT 'ユーザーID',
    is_owner BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'グループ作成者かどうか',
    ready_status BOOLEAN NOT NULL DEFAULT FALSE COMMENT '準備完了状態',
    ready_at TIMESTAMP NULL COMMENT '準備完了時刻',
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '参加時刻',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    -- 複合主キー
    PRIMARY KEY (group_id, user_id),

    -- インデックス
    INDEX idx_id (id),
    INDEX idx_group_id (group_id),
    INDEX idx_user_id (user_id),
    INDEX idx_is_owner (is_owner),
    INDEX idx_ready_status (ready_status),

    -- 外部キー制約
    CONSTRAINT fk_group_members_group_id
        FOREIGN KEY (group_id)
        REFERENCES `groups`(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_group_members_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='グループメンバーテーブル';
