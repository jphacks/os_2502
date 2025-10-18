-- collage_templatesテーブルの作成

CREATE TABLE IF NOT EXISTS collages_template (
    template_id CHAR(36) PRIMARY KEY COMMENT 'テンプレートID (UUID)',
    name VARCHAR(100) NOT NULL COMMENT 'テンプレート名',
    file_path VARCHAR(255) NOT NULL COMMENT 'テンプレートファイルパス',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    -- インデックス
    INDEX idx_name (name),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='コラージュテンプレートテーブル';
