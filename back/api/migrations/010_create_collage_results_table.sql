-- collage_resultsテーブルの作成

CREATE TABLE IF NOT EXISTS collage_results (
    result_id CHAR(36) PRIMARY KEY COMMENT '結果ID (UUID)',
    template_id CHAR(36) NOT NULL COMMENT 'テンプレートID',
    group_id CHAR(36) NOT NULL COMMENT 'グループID',
    file_url VARCHAR(500) NOT NULL COMMENT 'コラージュ画像URL',
    target_user_number INT NOT NULL COMMENT '対象ユーザー数',
    is_notification BOOLEAN NOT NULL DEFAULT FALSE COMMENT '通知済みフラグ',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',

    -- インデックス
    INDEX idx_template_id (template_id),
    INDEX idx_group_id (group_id),
    INDEX idx_created_at (created_at),
    INDEX idx_is_notification (is_notification),

    -- 外部キー制約
    CONSTRAINT fk_collage_results_template_id
        FOREIGN KEY (template_id)
        REFERENCES collages_template(template_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_collage_results_group_id
        FOREIGN KEY (group_id)
        REFERENCES `groups`(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='コラージュ結果テーブル';
