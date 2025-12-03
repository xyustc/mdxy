"""笔记相关路由"""
from fastapi import APIRouter, HTTPException, Query
from services.file_service import get_note_tree, get_note_content, search_notes

router = APIRouter(prefix="/api/notes", tags=["notes"])


@router.get("")
async def list_notes():
    """获取笔记目录树"""
    tree = get_note_tree()
    return {
        "success": True,
        "data": tree
    }


@router.get("/search")
async def search(q: str = Query(..., min_length=1, description="搜索关键词")):
    """搜索笔记"""
    results = search_notes(q)
    return {
        "success": True,
        "data": results,
        "total": len(results)
    }


@router.get("/{note_path:path}")
async def read_note(note_path: str):
    """获取指定笔记的内容"""
    content = get_note_content(note_path)
    
    if content is None:
        raise HTTPException(status_code=404, detail="笔记不存在")
    
    return {
        "success": True,
        "data": {
            "path": note_path,
            "content": content
        }
    }
