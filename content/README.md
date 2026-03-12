# 内容目录

此目录存放所有内容文件，通过 Git 管理。

## 目录结构

- `notes/` - 八股笔记（Markdown 文件）
- `articles/` - 技术博客文章（带 Front-matter 的 Markdown）

## 文章 Front-matter 格式

```yaml
---
title: "文章标题"
slug: "article-slug"
date: 2024-12-01
category: "分类"
tags: ["标签1", "标签2"]
summary: "文章摘要"
cover: "/images/cover.png"
status: published
---

文章正文...
```

## 注意事项

- 笔记文件直接通过文件系统读取
- 文章需要通过后端同步接口将元数据导入数据库
- 所有内容文件使用 UTF-8 编码
