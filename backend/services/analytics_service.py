"""访问统计服务"""
from sqlalchemy.orm import Session
from sqlalchemy import func, desc
from datetime import datetime, timedelta
from typing import List, Dict, Optional
from models.analytics import AccessLog


def create_access_log(db: Session, log_data: dict) -> AccessLog:
    """创建访问日志"""
    log = AccessLog(**log_data)
    db.add(log)
    db.commit()
    db.refresh(log)
    return log


def get_access_logs(
    db: Session,
    page: int = 1,
    limit: int = 50,
    start_date: Optional[str] = None,
    end_date: Optional[str] = None,
    ip: Optional[str] = None,
    path: Optional[str] = None
) -> Dict:
    """
    获取访问日志列表
    :param db: 数据库会话
    :param page: 页码
    :param limit: 每页数量
    :param start_date: 开始日期
    :param end_date: 结束日期
    :param ip: IP地址过滤
    :param path: 路径过滤
    """
    query = db.query(AccessLog)
    
    # 日期过滤
    if start_date:
        try:
            start = datetime.fromisoformat(start_date)
            query = query.filter(AccessLog.created_at >= start)
        except:
            pass
    
    if end_date:
        try:
            end = datetime.fromisoformat(end_date)
            query = query.filter(AccessLog.created_at <= end)
        except:
            pass
    
    # IP过滤
    if ip:
        query = query.filter(AccessLog.ip_address.like(f"%{ip}%"))
    
    # 路径过滤
    if path:
        query = query.filter(AccessLog.path.like(f"%{path}%"))
    
    # 总数
    total = query.count()
    
    # 分页
    offset = (page - 1) * limit
    logs = query.order_by(desc(AccessLog.created_at)).offset(offset).limit(limit).all()
    
    return {
        "logs": [log.to_dict() for log in logs],
        "total": total,
        "page": page,
        "limit": limit,
        "pages": (total + limit - 1) // limit
    }


def get_overview_stats(
    db: Session,
    start_date: Optional[str] = None,
    end_date: Optional[str] = None
) -> Dict:
    """
    获取统计概览
    """
    query = db.query(AccessLog)
    
    # 日期过滤
    if start_date:
        try:
            start = datetime.fromisoformat(start_date)
            query = query.filter(AccessLog.created_at >= start)
        except:
            pass
    
    if end_date:
        try:
            end = datetime.fromisoformat(end_date)
            query = query.filter(AccessLog.created_at <= end)
        except:
            pass
    
    # 总访问量
    total_visits = query.count()
    
    # 独立访客数（唯一IP）
    unique_visitors = db.query(func.count(func.distinct(AccessLog.ip_address))).filter(
        AccessLog.created_at >= datetime.fromisoformat(start_date) if start_date else True,
        AccessLog.created_at <= datetime.fromisoformat(end_date) if end_date else True
    ).scalar() or 0
    
    # 平均响应时间
    avg_response_time = db.query(func.avg(AccessLog.response_time)).filter(
        AccessLog.created_at >= datetime.fromisoformat(start_date) if start_date else True,
        AccessLog.created_at <= datetime.fromisoformat(end_date) if end_date else True
    ).scalar() or 0
    
    # Top访问页面
    top_pages = db.query(
        AccessLog.path,
        func.count(AccessLog.id).label('count')
    ).filter(
        AccessLog.created_at >= datetime.fromisoformat(start_date) if start_date else True,
        AccessLog.created_at <= datetime.fromisoformat(end_date) if end_date else True
    ).group_by(AccessLog.path).order_by(desc('count')).limit(10).all()
    
    # 访问趋势（按天）
    visitor_trends = db.query(
        func.date(AccessLog.created_at).label('date'),
        func.count(AccessLog.id).label('count')
    ).filter(
        AccessLog.created_at >= datetime.fromisoformat(start_date) if start_date else True,
        AccessLog.created_at <= datetime.fromisoformat(end_date) if end_date else True
    ).group_by(func.date(AccessLog.created_at)).order_by('date').all()
    
    # 设备统计
    device_stats = db.query(
        AccessLog.device_type,
        func.count(AccessLog.id).label('count')
    ).filter(
        AccessLog.created_at >= datetime.fromisoformat(start_date) if start_date else True,
        AccessLog.created_at <= datetime.fromisoformat(end_date) if end_date else True
    ).group_by(AccessLog.device_type).all()
    
    # 操作系统统计
    os_stats = db.query(
        AccessLog.os,
        func.count(AccessLog.id).label('count')
    ).filter(
        AccessLog.created_at >= datetime.fromisoformat(start_date) if start_date else True,
        AccessLog.created_at <= datetime.fromisoformat(end_date) if end_date else True
    ).group_by(AccessLog.os).order_by(desc('count')).limit(10).all()
    
    # 浏览器统计
    browser_stats = db.query(
        AccessLog.browser,
        func.count(AccessLog.id).label('count')
    ).filter(
        AccessLog.created_at >= datetime.fromisoformat(start_date) if start_date else True,
        AccessLog.created_at <= datetime.fromisoformat(end_date) if end_date else True
    ).group_by(AccessLog.browser).order_by(desc('count')).limit(10).all()
    
    return {
        "total_visits": total_visits,
        "unique_visitors": unique_visitors,
        "avg_response_time": round(avg_response_time, 2) if avg_response_time else 0,
        "top_pages": [{"path": p, "count": c} for p, c in top_pages],
        "visitor_trends": [{"date": str(d), "count": c} for d, c in visitor_trends],
        "device_stats": {d: c for d, c in device_stats},
        "os_stats": [{"os": o, "count": c} for o, c in os_stats],
        "browser_stats": [{"browser": b, "count": c} for b, c in browser_stats]
    }
