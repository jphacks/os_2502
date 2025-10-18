-- group_part_assignmentsテーブルの作成
-- グループメンバーへのパーツ割り当て

CREATE TABLE IF NOT EXISTS group_part_assignments (
    assignment_id CHAR(36) PRIMARY KEY COMMENT '割り当てID (UUID)',
    group_id CHAR(36) NOT NULL COMMENT 'グループID',
    user_id CHAR(36) NOT NULL COMMENT 'ユーザーID',
    part_id CHAR(36) NOT NULL COMMENT 'パーツID',
    collage_day DATE NOT NULL COMMENT 'コラージュ対象日',
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '割り当て日時',

    -- インデックス
    INDEX idx_group_id (group_id),
    INDEX idx_user_id (user_id),
    INDEX idx_part_id (part_id),
    INDEX idx_group_collage_day (group_id, collage_day),
    INDEX idx_user_group_day (user_id, group_id, collage_day),

    -- ユニーク制約（同じグループ・日付で同じパーツを複数人に割り当てない）
    UNIQUE KEY unique_group_day_part (group_id, collage_day, part_id),

    -- 外部キー制約
    CONSTRAINT fk_gpa_group_id
        FOREIGN KEY (group_id)
        REFERENCES `groups`(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_gpa_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_gpa_part_id
        FOREIGN KEY (part_id)
        REFERENCES template_parts(part_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='グループパーツ割り当てテーブル';
