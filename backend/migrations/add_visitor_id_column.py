"""
数据库迁移脚本：为 access_logs 表添加 visitor_id 字段
"""

from sqlalchemy import create_engine, MetaData, Table, Column, String
from sqlalchemy.dialects.sqlite import dialect
from config import DATABASE_URL

def upgrade():
    """添加 visitor_id 字段"""
    engine = create_engine(DATABASE_URL)
    metadata = MetaData()
    
    # 反射现有表结构
    metadata.reflect(bind=engine)
    
    # 获取 access_logs 表
    access_logs_table = metadata.tables['access_logs']
    
    # 添加 visitor_id 列（如果尚不存在）
    if 'visitor_id' not in access_logs_table.c:
        with engine.connect() as conn:
            # SQLite 不支持直接添加列，需要重建表
            # 这里我们使用 ALTER TABLE 语句（适用于支持的数据库）
            try:
                conn.execute('ALTER TABLE access_logs ADD COLUMN visitor_id VARCHAR(36)')
                conn.execute('CREATE INDEX IF NOT EXISTS ix_access_logs_visitor_id ON access_logs (visitor_id)')
                conn.commit()
                print("Successfully added visitor_id column to access_logs table")
            except Exception as e:
                print(f"Error adding visitor_id column: {e}")
                # 对于 SQLite，可能需要更复杂的处理
                print("Note: For SQLite, you may need to manually add the column")

def downgrade():
    """降级操作（留空，因为我们不支持降级）"""
    pass

if __name__ == "__main__":
    upgrade()