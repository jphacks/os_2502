-- device_tokensテーブルの作成

CREATE TABLE IF NOT EXISTS device_tokens (
    id CHAR(36) PRIMARY KEY COMMENT 'デバイストークンID (UUID)',
    user_id CHAR(36) NOT NULL COMMENT 'ユーザーID',
    device_token VARCHAR(255) NOT NULL COMMENT 'デバイストークン',
    device_type ENUM('ios', 'android') NOT NULL COMMENT 'デバイスタイプ',
    device_name VARCHAR(100) NULL COMMENT 'デバイス名',
    is_active BOOLEAN NOT NULL DEFAULT TRUE COMMENT 'アクティブフラグ',
    last_used_at TIMESTAMP NULL COMMENT '最終使用日時',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    -- インデックス
    INDEX idx_user_id (user_id),
    INDEX idx_device_token (device_token),
    INDEX idx_is_active (is_active),
    INDEX idx_last_used_at (last_used_at),

    -- 外部キー制約
    CONSTRAINT fk_device_tokens_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='デバイストークンテーブル';
