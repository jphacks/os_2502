-- Add scheduled_capture_time column to groups table

ALTER TABLE `groups`
ADD COLUMN `scheduled_capture_time` TIMESTAMP NULL COMMENT '予定撮影時刻（全クライアント同期用）' AFTER `countdown_started_at`;

-- Add index for scheduled_capture_time
ALTER TABLE `groups`
ADD INDEX `idx_scheduled_capture_time` (`scheduled_capture_time`);
