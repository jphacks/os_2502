-- usersテーブルの作成

CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY COMMENT 'ユーザーID (UUID)',
    firebase_uid VARCHAR(128) NOT NULL UNIQUE COMMENT 'Firebase UID',
    name VARCHAR(15) NOT NULL COMMENT 'ユーザー名（表示名）',
    username VARCHAR(30) NULL UNIQUE COMMENT 'ユーザー名（公開ID、ユニーク）',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    -- インデックス
    INDEX idx_firebase_uid (firebase_uid),
    INDEX idx_username (username),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='ユーザーテーブル';
