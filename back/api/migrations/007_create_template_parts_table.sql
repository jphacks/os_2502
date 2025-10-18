-- template_partsテーブルの作成
-- テンプレートごとのコラージュパーツ定義

CREATE TABLE IF NOT EXISTS template_parts (
    part_id CHAR(36) PRIMARY KEY COMMENT 'パーツID (UUID)',
    template_id CHAR(36) NOT NULL COMMENT 'テンプレートID',
    part_number INT NOT NULL COMMENT 'パーツ番号（1から始まる連番）',
    part_name VARCHAR(100) NULL COMMENT 'パーツ名（例: "左上", "中央"）',
    position_x INT NOT NULL COMMENT 'X座標',
    position_y INT NOT NULL COMMENT 'Y座標',
    width INT NOT NULL COMMENT '幅',
    height INT NOT NULL COMMENT '高さ',
    description TEXT NULL COMMENT 'パーツの説明（どんな写真を撮るか等）',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    -- インデックス
    INDEX idx_template_id (template_id),
    INDEX idx_template_part_number (template_id, part_number),

    -- ユニーク制約（同じテンプレート内でパーツ番号は一意）
    UNIQUE KEY unique_template_part_number (template_id, part_number),

    -- 外部キー制約
    CONSTRAINT fk_template_parts_template_id
        FOREIGN KEY (template_id)
        REFERENCES collages_template(template_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='テンプレートパーツテーブル';
