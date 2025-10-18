-- upload_images_collage_resultsテーブルの作成
-- 画像とコラージュ結果の関連テーブル

CREATE TABLE IF NOT EXISTS upload_images_collage_results (
    image_id CHAR(36) NOT NULL COMMENT '画像ID',
    result_id CHAR(36) NOT NULL COMMENT '結果ID',
    position_x INT NOT NULL COMMENT 'X座標',
    position_y INT NOT NULL COMMENT 'Y座標',
    width INT NOT NULL COMMENT '幅',
    height INT NOT NULL COMMENT '高さ',
    sort_order INT NOT NULL COMMENT '表示順序',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',

    -- 複合主キー
    PRIMARY KEY (image_id, result_id),

    -- インデックス
    INDEX idx_result_id (result_id),
    INDEX idx_sort_order (sort_order),

    -- 外部キー制約
    CONSTRAINT fk_uicr_image_id
        FOREIGN KEY (image_id)
        REFERENCES upload_images(image_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_uicr_result_id
        FOREIGN KEY (result_id)
        REFERENCES collage_results(result_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='画像とコラージュ結果の関連テーブル';
