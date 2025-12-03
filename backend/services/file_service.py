"""文件服务 - 处理笔记文件的读取"""
import os
from pathlib import Path
from typing import Optional
from config import NOTES_DIR, ALLOWED_EXTENSIONS


def get_note_tree() -> list:
    """
    获取笔记目录树
    返回结构: [{ name, type, path, children? }]
    """
    return _scan_directory(NOTES_DIR)


def _scan_directory(directory: Path, relative_base: Path = None) -> list:
    """递归扫描目录"""
    if relative_base is None:
        relative_base = NOTES_DIR
    
    items = []
    
    if not directory.exists():
        return items
    
    # 获取目录内容并排序（文件夹优先，然后按名称排序）
    entries = sorted(
        directory.iterdir(),
        key=lambda x: (not x.is_dir(), x.name.lower())
    )
    
    for entry in entries:
        # 跳过隐藏文件和目录
        if entry.name.startswith('.'):
            continue
        
        relative_path = entry.relative_to(relative_base)
        
        if entry.is_dir():
            children = _scan_directory(entry, relative_base)
            # 只添加非空文件夹或包含md文件的文件夹
            if children:
                items.append({
                    "name": entry.name,
                    "type": "directory",
                    "path": str(relative_path).replace("\\", "/"),
                    "children": children
                })
        elif entry.suffix.lower() in ALLOWED_EXTENSIONS:
            items.append({
                "name": entry.stem,  # 不带扩展名的文件名
                "type": "file",
                "path": str(relative_path).replace("\\", "/"),
            })
    
    return items


def get_note_content(note_path: str) -> Optional[str]:
    """
    获取笔记内容
    :param note_path: 相对于 notes 目录的路径
    :return: 笔记内容或 None
    """
    # 安全检查：防止目录穿越攻击
    safe_path = _safe_path(note_path)
    if safe_path is None:
        return None
    
    if not safe_path.exists():
        return None
    
    if safe_path.suffix.lower() not in ALLOWED_EXTENSIONS:
        return None
    
    try:
        return safe_path.read_text(encoding='utf-8')
    except Exception:
        return None


def _safe_path(note_path: str) -> Optional[Path]:
    """
    验证并返回安全的文件路径
    防止目录穿越攻击
    """
    try:
        # 规范化路径
        normalized = Path(note_path).as_posix()
        # 移除开头的斜杠
        normalized = normalized.lstrip('/')
        
        # 构建完整路径
        full_path = (NOTES_DIR / normalized).resolve()
        
        # 确保路径在 NOTES_DIR 内
        if not str(full_path).startswith(str(NOTES_DIR.resolve())):
            return None
        
        return full_path
    except Exception:
        return None


def search_notes(keyword: str) -> list:
    """
    搜索笔记内容
    :param keyword: 搜索关键词
    :return: 匹配的笔记列表
    """
    results = []
    keyword_lower = keyword.lower()
    
    for md_file in NOTES_DIR.rglob("*.md"):
        try:
            content = md_file.read_text(encoding='utf-8')
            if keyword_lower in content.lower() or keyword_lower in md_file.stem.lower():
                relative_path = md_file.relative_to(NOTES_DIR)
                # 提取匹配的上下文
                context = _extract_context(content, keyword)
                results.append({
                    "name": md_file.stem,
                    "path": str(relative_path).replace("\\", "/"),
                    "context": context
                })
        except Exception:
            continue
    
    return results


def _extract_context(content: str, keyword: str, context_length: int = 100) -> str:
    """提取关键词周围的上下文"""
    keyword_lower = keyword.lower()
    content_lower = content.lower()
    
    pos = content_lower.find(keyword_lower)
    if pos == -1:
        # 如果内容中没找到，返回开头部分
        return content[:context_length] + "..." if len(content) > context_length else content
    
    start = max(0, pos - context_length // 2)
    end = min(len(content), pos + len(keyword) + context_length // 2)
    
    context = content[start:end]
    if start > 0:
        context = "..." + context
    if end < len(content):
        context = context + "..."
    
    return context
