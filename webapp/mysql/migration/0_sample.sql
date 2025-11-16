-- このファイルに記述されたSQLコマンドが、マイグレーション時に実行されます。
-- -- このファイルに記述されたSQLコマンドが、マイグレーション時に実行されます。
-- -- -------------------------
-- -- orders 表索引
-- -- -------------------------
-- -- 加速 user_id 过滤和 JOIN
-- CREATE INDEX idx_orders_user_product ON orders(user_id, product_id);

-- -- 加速按 created_at 排序
-- CREATE INDEX idx_orders_user_created_at ON orders(user_id, created_at);

-- -- 加速按 shipped_status 排序
-- CREATE INDEX idx_orders_user_shipped_status ON orders(user_id, shipped_status);

-- -- 加速按 arrived_at 排序
-- CREATE INDEX idx_orders_user_arrived_at ON orders(user_id, arrived_at);

-- -- -------------------------
-- -- products 表索引
-- -- -------------------------
-- -- 加速 product_id JOIN（如果不是主键）
-- CREATE INDEX idx_products_product_id ON products(product_id);

-- -- 加速名称搜索（前缀 LIKE 'abc%' 可用）
-- CREATE INDEX idx_products_name ON products(name);

-- -- 如果有包含搜索 LIKE '%abc%' 需求，添加全文索引
-- -- ALTER TABLE products ADD FULLTEXT INDEX ft_name (name);

-- products 表
CREATE INDEX idx_name_id ON products(name, product_id);

-- orders 表
CREATE INDEX idx_orders_user_created_at ON orders(user_id, created_at);
CREATE INDEX idx_orders_user_shipped_status ON orders(user_id, shipped_status);
CREATE INDEX idx_orders_user_arrived_at ON orders(user_id, arrived_at);


-- session_id + expires_at 快速查询有效会话
CREATE INDEX idx_user_sessions_session_expires ON user_sessions(session_uuid, expires_at);

-- user_id + expires_at 查询某用户有效会话
CREATE INDEX idx_user_sessions_user_expires ON user_sessions(user_id, expires_at);

-- -------------------------
-- users 表
-- -------------------------
-- user_name 唯一索引/快速查找
CREATE UNIQUE INDEX idx_users_user_name ON users(user_name);
