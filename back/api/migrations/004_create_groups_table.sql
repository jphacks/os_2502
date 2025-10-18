-- groupsテーブルの作成

CREATE TABLE IF NOT EXISTS `groups` (
    id CHAR(36) PRIMARY KEY COMMENT 'グループID (UUID)',
    owner_user_id CHAR(36) NOT NULL COMMENT 'グループ作成者のユーザーID',
    name VARCHAR(15) NOT NULL COMMENT 'グループ名',
    group_type ENUM('local_temporary', 'global_temporary', 'permanent') NOT NULL DEFAULT 'global_temporary' COMMENT 'グループタイプ',
    status ENUM('recruiting', 'ready_check', 'countdown', 'photo_taking', 'completed', 'expired') NOT NULL DEFAULT 'recruiting' COMMENT 'グループステータス',
    max_member INT NOT NULL COMMENT '最大メンバー数',
    current_member_count INT NOT NULL DEFAULT 0 COMMENT '現在のメンバー数',
    invitation_token CHAR(36) NOT NULL UNIQUE COMMENT '招待リンク用トークン',
    finalized_at TIMESTAMP NULL COMMENT 'メンバー確定時刻',
    countdown_started_at TIMESTAMP NULL COMMENT 'カウントダウン開始時刻',
    expires_at TIMESTAMP NULL COMMENT '有効期限（一時グループ用）',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    -- インデックス
    INDEX idx_owner_user_id (owner_user_id),
    INDEX idx_status (status),
    INDEX idx_group_type (group_type),
    INDEX idx_invitation_token (invitation_token),
    INDEX idx_created_at (created_at),

    -- 外部キー制約
    CONSTRAINT fk_groups_owner_user_id
        FOREIGN KEY (owner_user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='グループテーブル';
