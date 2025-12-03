"""
Markdown 笔记系统 - 后端入口
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from config import CORS_ORIGINS, NOTES_DIR
from routers import notes

# 创建 FastAPI 应用
app = FastAPI(
    title="Markdown 笔记系统",
    description="在线 Markdown 笔记查看系统 API",
    version="1.0.0"
)

# 配置 CORS 跨域
app.add_middleware(
    CORSMiddleware,
    allow_origins=CORS_ORIGINS,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 注册路由
app.include_router(notes.router)


@app.get("/")
async def root():
    """API 根路径"""
    return {
        "message": "Markdown 笔记系统 API",
        "docs": "/docs",
        "notes_dir": str(NOTES_DIR)
    }


@app.get("/health")
async def health_check():
    """健康检查"""
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
