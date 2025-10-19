-- Add template_id column to groups table for synchronized template selection
ALTER TABLE `groups` ADD COLUMN `template_id` CHAR(36) NULL COMMENT '選択されたテンプレートID（グループ内で統一）' AFTER `scheduled_capture_time`;
