"""访问日志数据模型"""
from sqlalchemy import Column, Integer, String, Float, DateTime, Index
from sqlalchemy.ext.declarative import declarative_base
from datetime import datetime

Base = declarative_base()


class AccessLog(Base):
    """访问日志表"""
    __tablename__ = 'access_logs'
    
    id = Column(Integer, primary_key=True, autoincrement=True)
    ip_address = Column(String(45), index=True)  # 支持IPv6
    user_agent = Column(String(500))
    path = Column(String(500), index=True)
    method = Column(String(10))
    status_code = Column(Integer)
    response_time = Column(Float)  # 毫秒
    referer = Column(String(500))
    device_type = Column(String(50))  # PC/Mobile/Tablet
    os = Column(String(50))
    browser = Column(String(50))
    created_at = Column(DateTime, default=datetime.utcnow, index=True)
    
    def to_dict(self):
        """转换为字典"""
        return {
            "id": self.id,
            "ip_address": self.ip_address,
            "user_agent": self.user_agent,
            "path": self.path,
            "method": self.method,
            "status_code": self.status_code,
            "response_time": round(self.response_time, 2) if self.response_time else 0,
            "referer": self.referer,
            "device_type": self.device_type,
            "os": self.os,
            "browser": self.browser,
            "created_at": self.created_at.isoformat() if self.created_at else None
        }
