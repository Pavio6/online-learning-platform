-- 移除 comments 表的 parent_comment_id 外键约束
-- 因为评论是跨分片的，父评论可能在另一个分支数据库中
-- 需要在每个分支节点数据库中执行

-- 查找并删除外键约束
-- PostgreSQL 会自动生成约束名称，通常是 "comments_parent_comment_id_fkey"
DO $$
DECLARE
    constraint_name TEXT;
BEGIN
    -- 查找外键约束名称
    SELECT constraint_name INTO constraint_name
    FROM information_schema.table_constraints
    WHERE table_schema = 'public'
      AND table_name = 'comments'
      AND constraint_type = 'FOREIGN KEY'
      AND constraint_name LIKE '%parent_comment_id%';
    
    -- 如果找到约束，则删除
    IF constraint_name IS NOT NULL THEN
        EXECUTE 'ALTER TABLE comments DROP CONSTRAINT IF EXISTS ' || constraint_name;
        RAISE NOTICE 'Dropped constraint: %', constraint_name;
    ELSE
        RAISE NOTICE 'No foreign key constraint found for parent_comment_id';
    END IF;
END $$;

-- 备用方法：直接尝试删除常见的约束名称
ALTER TABLE comments DROP CONSTRAINT IF EXISTS comments_parent_comment_id_fkey;
ALTER TABLE comments DROP CONSTRAINT IF EXISTS comments_parent_comment_id_fk;

