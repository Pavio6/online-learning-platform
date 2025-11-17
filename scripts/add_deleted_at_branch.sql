-- 为分支节点所有表添加deleted_at字段
-- 在每个分支节点数据库中执行（learning_branch1, learning_branch2等）

-- branches表
ALTER TABLE branches ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_branches_deleted_at ON branches(deleted_at);

-- users表
ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- answers表
ALTER TABLE answers ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_answers_deleted_at ON answers(deleted_at);

-- comments表
ALTER TABLE comments ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_comments_deleted_at ON comments(deleted_at);

-- learning表
ALTER TABLE learning ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_learning_deleted_at ON learning(deleted_at);

-- courses表（分支节点的副本）
ALTER TABLE courses ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_courses_deleted_at ON courses(deleted_at);

-- chapters表（分支节点的副本）
ALTER TABLE chapters ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_chapters_deleted_at ON chapters(deleted_at);

-- lessons表（分支节点的副本）
ALTER TABLE lessons ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_lessons_deleted_at ON lessons(deleted_at);

-- tasks表（分支节点的副本）
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_tasks_deleted_at ON tasks(deleted_at);

