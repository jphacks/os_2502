-- upload_imagesテーブルの作成

CREATE TABLE IF NOT EXISTS upload_images (
    image_id CHAR(36) PRIMARY KEY COMMENT '画像ID (UUID)',
    file_url VARCHAR(500) NOT NULL COMMENT '画像ファイルURL',
    group_id CHAR(36) NOT NULL COMMENT 'グループID',
    user_id CHAR(36) NOT NULL COMMENT 'アップロードユーザーID',
    part_id CHAR(36) NULL COMMENT 'パーツID（どのパーツ用の画像か）',
    collage_day DATE NOT NULL COMMENT 'コラージュ対象日',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'アップロード日時',

    -- インデックス
    INDEX idx_group_id (group_id),
    INDEX idx_user_id (user_id),
    INDEX idx_part_id (part_id),
    INDEX idx_collage_day (collage_day),
    INDEX idx_created_at (created_at),
    INDEX idx_group_collage_day (group_id, collage_day),

    -- 外部キー制約
    CONSTRAINT fk_upload_images_group_id
        FOREIGN KEY (group_id)
        REFERENCES `groups`(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_upload_images_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_upload_images_part_id
        FOREIGN KEY (part_id)
        REFERENCES template_parts(part_id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='アップロード画像テーブル';
