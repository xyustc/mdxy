"""访问统计路由"""
from fastapi import APIRouter, Depends, Query
from sqlalchemy.orm import Session
from typing import Optional
from datetime import datetime, timedelta, timezone

from database import get_db
from utils.jwt_utils import verify_admin_dependency
from services.analytics_service import get_access_logs, get_overview_stats


router = APIRouter(prefix="/api/admin/analytics", tags=["analytics"])


@router.get("/logs")
async def list_logs(
    page: int = Query(1, ge=1, description="页码"),
    limit: int = Query(50, ge=1, le=100, description="每页数量"),
    start_date: Optional[str] = Query(None, description="开始日期 ISO格式"),
    end_date: Optional[str] = Query(None, description="结束日期 ISO格式"),
    ip: Optional[str] = Query(None, description="IP地址过滤"),
    path: Optional[str] = Query(None, description="路径过滤"),
    db: Session = Depends(get_db),
    _: dict = Depends(verify_admin_dependency)
):
    """获取访问日志列表"""
    result = get_access_logs(
        db=db,
        page=page,
        limit=limit,
        start_date=start_date,
        end_date=end_date,
        ip=ip,
        path=path
    )
    
    return {
        "success": True,
        "data": result["logs"],
        "total": result["total"],
        "page": result["page"],
        "limit": result["limit"],
        "pages": result["pages"]
    }


@router.get("/overview")
async def overview(
    start_date: Optional[str] = Query(None, description="开始日期 ISO格式"),
    end_date: Optional[str] = Query(None, description="结束日期 ISO格式"),
    db: Session = Depends(get_db),
    _: dict = Depends(verify_admin_dependency)
):
    """获取统计概览"""
    # 如果没有指定日期，默认最近30天
    if not start_date:
        start_date = datetime.now(timezone.utc) - timedelta(days=30)
        start_date = start_date.isoformat()
    if not end_date:
        end_date = datetime.now(timezone.utc).isoformat()
    
    stats = get_overview_stats(db=db, start_date=start_date, end_date=end_date)
    
    return {
        "success": True,
        "data": stats
    }
