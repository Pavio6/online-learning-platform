-- 初始化branches表数据
-- 需要在每个分支节点的数据库中执行

-- 分支节点1
INSERT INTO branches (branch_id, branch_name, created_at, updated_at) 
VALUES (1, '北京校区', NOW(), NOW())
ON CONFLICT (branch_id) DO NOTHING;

-- 分支节点2
INSERT INTO branches (branch_id, branch_name, created_at, updated_at) 
VALUES (2, '上海校区', NOW(), NOW())
ON CONFLICT (branch_id) DO NOTHING;

