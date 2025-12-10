"""  
Markdown 笔记系统 - 后端入口
"""
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from config import CORS_ORIGINS, NOTES_DIR
from routers import notes, admin, analytics
from middleware import AccessLogMiddleware
from database import init_database
import os
import time

# 创建 FastAPI 应用
app = FastAPI(
    title="Xingyu的笔记",
    description="Xingyu的在线 Markdown 笔记查看系统 API",
    version="1.1.0"
)

# 初始化数据库
init_database()

# 配置 CORS 跨域
app.add_middleware(
    CORSMiddleware,
    allow_origins=CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 添加访问日志中间件
app.add_middleware(AccessLogMiddleware)

# 注册路由
app.include_router(notes.router)
app.include_router(admin.router)
app.include_router(analytics.router)


@app.get("/")
async def root():
    """API 根路径"""
    return {
        "message": "Markdown 笔记系统 API",
        "docs": "/docs",
        "notes_dir": str(NOTES_DIR),
        "timestamp": time.time()
    }


@app.get("/health")
async def health_check():
    """增强的健康检查"""
    try:
        # 检查笔记目录是否存在且可访问
        if not NOTES_DIR.exists():
            raise HTTPException(status_code=503, detail="笔记目录不存在")
        
        if not os.access(NOTES_DIR, os.R_OK):
            raise HTTPException(status_code=503, detail="笔记目录不可读")
            
        return {
            "status": "ok",
            "timestamp": time.time(),
            "notes_dir_accessible": True,
            "service": "mdxy-backend"
        }
    except Exception as e:
        raise HTTPException(status_code=503, detail=f"健康检查失败: {str(e)}")


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)