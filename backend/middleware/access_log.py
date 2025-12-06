"""访问日志中间件"""
import time
from fastapi import Request
from starlette.middleware.base import BaseHTTPMiddleware
from database import SessionLocal
from utils.parser import parse_user_agent
from services.analytics_service import create_access_log


class AccessLogMiddleware(BaseHTTPMiddleware):
    """访问日志记录中间件"""
    
    async def dispatch(self, request: Request, call_next):
        # 排除不需要记录的路径
        excluded_paths = [
            "/api/admin/",  # 管理接口
            "/docs",        # API文档
            "/openapi.json",
            "/health",
            "/favicon.ico"
        ]
        
        # 检查是否需要记录
        should_log = True
        for path in excluded_paths:
            if request.url.path.startswith(path):
                should_log = False
                break
        
        # 记录开始时间
        start_time = time.time()
        
        # 执行请求
        response = await call_next(request)
        
        # 计算响应时间
        response_time = (time.time() - start_time) * 1000  # 转换为毫秒
        
        # 记录访问日志（仅记录笔记访问）
        if should_log and request.url.path.startswith("/api/notes"):
            try:
                # 解析User-Agent
                user_agent_string = request.headers.get("user-agent", "")
                ua_info = parse_user_agent(user_agent_string)
                
                # 获取客户端IP
                client_ip = request.client.host if request.client else "unknown"
                
                # 准备日志数据
                log_data = {
                    "ip_address": client_ip,
                    "user_agent": user_agent_string,
                    "path": request.url.path,
                    "method": request.method,
                    "status_code": response.status_code,
                    "response_time": response_time,
                    "referer": request.headers.get("referer", ""),
                    "device_type": ua_info["device_type"],
                    "os": ua_info["os"],
                    "browser": ua_info["browser"]
                }
                
                # 异步写入数据库
                db = SessionLocal()
                try:
                    create_access_log(db, log_data)
                finally:
                    db.close()
            except Exception as e:
                # 记录日志失败不影响正常请求
                print(f"访问日志记录失败: {str(e)}")
        
        return response
