"""管理员路由 - 认证和文件管理"""
from fastapi import APIRouter, HTTPException, Depends, UploadFile, File, Form
from pydantic import BaseModel
from typing import Optional
from pathlib import Path
import shutil

from services.auth_service import authenticate_admin
from services.file_service import _safe_path
from utils.jwt_utils import verify_admin_dependency
from config import NOTES_DIR, ALLOWED_EXTENSIONS, MAX_UPLOAD_SIZE


router = APIRouter(prefix="/api/admin", tags=["admin"])


# ========== 请求模型 ==========
class LoginRequest(BaseModel):
    password: str


class MkdirRequest(BaseModel):
    path: str


class MoveRequest(BaseModel):
    old_path: str
    new_path: str


# ========== 认证接口 ==========
@router.post("/login")
async def login(request: LoginRequest):
    """管理员登录"""
    result = authenticate_admin(request.password)
    
    if not result["success"]:
        raise HTTPException(status_code=401, detail=result["message"])
    
    return {
        "success": True,
        "data": {
            "token": result["token"]
        },
        "message": result["message"]
    }


@router.get("/verify")
async def verify_token(payload: dict = Depends(verify_admin_dependency)):
    """验证token是否有效"""
    return {
        "success": True,
        "data": payload,
        "message": "Token有效"
    }


# ========== 文件管理接口 ==========
@router.post("/notes/upload")
async def upload_note(
    file: UploadFile = File(...),
    path: str = Form(""),
    _: dict = Depends(verify_admin_dependency)
):
    """
    上传笔记文件
    :param file: 上传的文件
    :param path: 目标路径（相对于notes目录）
    """
    # 验证文件扩展名
    file_ext = Path(file.filename).suffix.lower()
    if file_ext not in ALLOWED_EXTENSIONS:
        raise HTTPException(
            status_code=400,
            detail=f"不支持的文件类型，仅允许: {', '.join(ALLOWED_EXTENSIONS)}"
        )
    
    # 构建目标路径
    if path:
        target_dir = _safe_path(path)
        if target_dir is None:
            raise HTTPException(status_code=400, detail="无效的路径")
        
        # 如果目标是文件路径，获取其父目录
        if target_dir.suffix:
            target_dir = target_dir.parent
    else:
        target_dir = NOTES_DIR
    
    # 确保目录存在
    target_dir.mkdir(parents=True, exist_ok=True)
    
    # 构建完整文件路径
    target_file = target_dir / file.filename
    
    # 检查文件是否已存在
    if target_file.exists():
        raise HTTPException(status_code=400, detail="文件已存在")
    
    # 读取文件内容并检查大小
    content = await file.read()
    if len(content) > MAX_UPLOAD_SIZE:
        raise HTTPException(
            status_code=400,
            detail=f"文件过大，最大允许 {MAX_UPLOAD_SIZE / 1024 / 1024}MB"
        )
    
    # 写入文件
    try:
        target_file.write_bytes(content)
        relative_path = target_file.relative_to(NOTES_DIR)
        
        return {
            "success": True,
            "data": {
                "path": str(relative_path).replace("\\", "/"),
                "name": file.filename,
                "size": len(content)
            },
            "message": "上传成功"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"上传失败: {str(e)}")


@router.delete("/notes/{note_path:path}")
async def delete_note(
    note_path: str,
    _: dict = Depends(verify_admin_dependency)
):
    """删除笔记文件"""
    safe_path = _safe_path(note_path)
    
    if safe_path is None:
        raise HTTPException(status_code=400, detail="无效的路径")
    
    if not safe_path.exists():
        raise HTTPException(status_code=404, detail="文件不存在")
    
    try:
        if safe_path.is_file():
            safe_path.unlink()
        elif safe_path.is_dir():
            shutil.rmtree(safe_path)
        
        return {
            "success": True,
            "message": "删除成功"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"删除失败: {str(e)}")


@router.post("/notes/mkdir")
async def create_directory(
    request: MkdirRequest,
    _: dict = Depends(verify_admin_dependency)
):
    """创建目录"""
    safe_path = _safe_path(request.path)
    
    if safe_path is None:
        raise HTTPException(status_code=400, detail="无效的路径")
    
    if safe_path.exists():
        raise HTTPException(status_code=400, detail="目录已存在")
    
    try:
        safe_path.mkdir(parents=True, exist_ok=True)
        relative_path = safe_path.relative_to(NOTES_DIR)
        
        return {
            "success": True,
            "data": {
                "path": str(relative_path).replace("\\", "/")
            },
            "message": "创建成功"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"创建失败: {str(e)}")


@router.put("/notes/move")
async def move_note(
    request: MoveRequest,
    _: dict = Depends(verify_admin_dependency)
):
    """移动/重命名文件"""
    old_path = _safe_path(request.old_path)
    new_path = _safe_path(request.new_path)
    
    if old_path is None or new_path is None:
        raise HTTPException(status_code=400, detail="无效的路径")
    
    if not old_path.exists():
        raise HTTPException(status_code=404, detail="源文件不存在")
    
    if new_path.exists():
        raise HTTPException(status_code=400, detail="目标路径已存在")
    
    try:
        # 确保目标目录存在
        new_path.parent.mkdir(parents=True, exist_ok=True)
        
        # 移动文件
        shutil.move(str(old_path), str(new_path))
        
        relative_path = new_path.relative_to(NOTES_DIR)
        
        return {
            "success": True,
            "data": {
                "path": str(relative_path).replace("\\", "/")
            },
            "message": "移动成功"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"移动失败: {str(e)}")
