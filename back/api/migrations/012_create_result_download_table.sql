-- result_downloadテーブルの作成
-- コラージュ結果のダウンロード履歴

CREATE TABLE IF NOT EXISTS result_download (
    result_id CHAR(36) NOT NULL COMMENT '結果ID',
    user_id CHAR(36) NOT NULL COMMENT 'ユーザーID',
    downloaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'ダウンロード日時',

    -- 複合主キー
    PRIMARY KEY (result_id, user_id),

    -- インデックス
    INDEX idx_result_id (result_id),
    INDEX idx_user_id (user_id),
    INDEX idx_downloaded_at (downloaded_at),

    -- 外部キー制約
    CONSTRAINT fk_result_download_result_id
        FOREIGN KEY (result_id)
        REFERENCES collage_results(result_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_result_download_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='コラージュ結果ダウンロード履歴テーブル';
