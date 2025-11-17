-- 为中央服务器所有表添加deleted_at字段
-- 在中央服务器数据库（learning_central）中执行

-- courses表
ALTER TABLE courses ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_courses_deleted_at ON courses(deleted_at);

-- chapters表
ALTER TABLE chapters ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_chapters_deleted_at ON chapters(deleted_at);

-- lessons表
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_lessons_deleted_at ON lessons(deleted_at);

-- tasks表
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);

