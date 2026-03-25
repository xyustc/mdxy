export interface CheatSheetRow {
  key: string
  desc?: string
  addedAt?: string
  keyVariant?: 'default' | 'paste-image'
}

export interface CheatSheetGroup {
  title: string
  rows: CheatSheetRow[]
}

export interface CheatSheetSection {
  id: string
  title: string
  theme: 'keyboard' | 'mcp' | 'slash' | 'memory' | 'workflows' | 'config' | 'skills' | 'cli'
  groups: CheatSheetGroup[]
}

export interface CheatSheetFooterRow {
  label: string
  items: Array<{ code: string; text?: string }>
}

export interface CheatSheetCodeItem {
  code: string
  text?: string
}

export interface CheatSheetMeta {
  sourceUrl: string
  sourceTitle: string
  version: string
  updatedAt: string
}

export interface CheatSheetContent {
  meta: CheatSheetMeta
  changelog: CheatSheetCodeItem[]
  footer: CheatSheetFooterRow[]
  columns: CheatSheetSection[][]
}

export const claudeCodeCheatSheetMeta = {
  sourceUrl: 'https://banwagong1.com/claude-code.html',
  sourceTitle: 'Claude Code 速查表',
  version: 'Claude Code v2.1.81',
  updatedAt: '2026年3月23日'
}

export const claudeCodeCheatSheetChangelog = [
  { code: '--bare', text: '标志，最小化无头模式（无 hooks/LSP/插件）' },
  { code: '--channels', text: '权限中继与 MCP 推送消息（预览）' },
  { code: 'effort', text: '前置元数据可用于技能和斜杠命令' },
  { code: '/fork', text: '重命名为 /branch（别名保留）' },
  { code: 'SendMessage', text: '自动恢复已停止的代理' }
]

export const claudeCodeCheatSheetFooter: CheatSheetFooterRow[] = [
  {
    label: '权限模式',
    items: [
      { code: 'default', text: '提示确认' },
      { code: 'acceptEdits', text: '自动接受编辑' },
      { code: 'plan', text: '只读' },
      { code: 'dontAsk', text: '未允许则拒绝' },
      { code: 'bypassPermissions', text: '跳过所有' }
    ]
  },
  {
    label: '重要环境变量',
    items: [
      { code: 'ANTHROPIC_API_KEY' },
      { code: 'ANTHROPIC_MODEL' },
      { code: 'CLAUDE_CODE_EFFORT_LEVEL', text: '(low/med/high)' },
      { code: 'MAX_THINKING_TOKENS', text: '(0=关闭)' },
      { code: 'CLAUDE_CODE_MAX_OUTPUT_TOKENS', text: '(默认 32K)' },
      { code: 'CLAUDE_CODE_DISABLE_CRON' }
    ]
  }
]

export const claudeCodeCheatSheetColumns: CheatSheetSection[][] = [
  [
    {
      id: 'keyboard',
      title: '⌨️ 键盘快捷键',
      theme: 'keyboard',
      groups: [
        {
          title: '常规控制',
          rows: [
            { key: '[Ctrl][C]', desc: '取消输入/生成' },
            { key: '[Ctrl][D]', desc: '退出会话' },
            { key: '[Ctrl][L]', desc: '清屏' },
            { key: '[Ctrl][O]', desc: '切换详细输出' },
            { key: '[Ctrl][R]', desc: '反向搜索历史' },
            { key: '[Ctrl][G]', desc: '在编辑器中打开提示' },
            { key: '[Ctrl][B]', desc: '后台运行任务' },
            { key: '[Ctrl][T]', desc: '切换任务列表' },
            { key: '[Ctrl][V]', desc: '粘贴图片', keyVariant: 'paste-image' },
            { key: '[Ctrl][F]', desc: '终止后台代理（按两次）' },
            { key: '[Esc][Esc]', desc: '回退 / 撤销' }
          ]
        },
        {
          title: '模式切换',
          rows: [
            { key: '[Shift][Tab]', desc: '循环切换权限模式' },
            { key: '[Alt][P]', desc: '切换模型' },
            { key: '[Alt][T]', desc: '切换思考模式' }
          ]
        },
        {
          title: '输入',
          rows: [
            { key: '[\\][Enter]', desc: '换行（快捷）' },
            { key: '[Ctrl][J]', desc: '换行（控制序列）' }
          ]
        },
        {
          title: '前缀',
          rows: [
            { key: '/', desc: '斜杠命令' },
            { key: '!', desc: '直接执行 bash' },
            { key: '@', desc: '文件引用 + 自动补全' }
          ]
        },
        {
          title: '会话选择器',
          rows: [
            { key: '[↑↓]', desc: '导航' },
            { key: '[←→]', desc: '展开/折叠' },
            { key: '[P]', desc: '预览' },
            { key: '[R]', desc: '重命名' },
            { key: '[/]', desc: '搜索' },
            { key: '[A]', desc: '所有项目' },
            { key: '[B]', desc: '当前分支' }
          ]
        }
      ]
    },
    {
      id: 'mcp',
      title: '🔌 MCP 服务器',
      theme: 'mcp',
      groups: [
        {
          title: '添加服务器',
          rows: [
            { key: '--transport http', desc: '远程 HTTP（推荐）' },
            { key: '--transport stdio', desc: '本地进程' },
            { key: '--transport sse', desc: '远程 SSE' }
          ]
        },
        {
          title: '作用域',
          rows: [
            { key: 'Local', desc: '~/.claude.json（每个项目）' },
            { key: 'Project', desc: '.mcp.json（共享 / VCS）' },
            { key: 'User', desc: '~/.claude.json（全局）' }
          ]
        },
        {
          title: '管理',
          rows: [
            { key: '/mcp', desc: '交互式界面' },
            { key: 'claude mcp list', desc: '列出所有服务器' },
            { key: 'claude mcp serve', desc: 'CC 作为 MCP 服务器' },
            { key: 'Elicitation', desc: '服务器在任务中请求输入', addedAt: '2026-03-23' }
          ]
        }
      ]
    }
  ],
  [
    {
      id: 'slash',
      title: '⚡ 斜杠命令',
      theme: 'slash',
      groups: [
        {
          title: '会话',
          rows: [
            { key: '/clear', desc: '清除对话' },
            { key: '/compact [focus]', desc: '压缩上下文' },
            { key: '/resume', desc: '恢复 / 切换会话' },
            { key: '/rename [name]', desc: '命名当前会话' },
            { key: '/branch [name]', desc: '分支对话（/fork 别名）' },
            { key: '/cost', desc: 'Token 用量统计' },
            { key: '/context', desc: '可视化上下文（网格）' },
            { key: '/diff', desc: '交互式差异查看器' },
            { key: '/copy', desc: '复制上次回复' },
            { key: '/export', desc: '导出对话' }
          ]
        },
        {
          title: '配置',
          rows: [
            { key: '/config', desc: '打开设置' },
            { key: '/model [model]', desc: '切换模型（←→ 努力度）' },
            { key: '/fast [on|off]', desc: '切换快速模式' },
            { key: '/vim', desc: '切换 vim 模式' },
            { key: '/theme', desc: '更改颜色主题' },
            { key: '/permissions', desc: '查看 / 更新权限' },
            { key: '/effort [level]', desc: '设置努力度（low/med/high）', addedAt: '2026-03-23' },
            { key: '/color [color]', desc: '设置提示栏颜色' }
          ]
        },
        {
          title: '工具',
          rows: [
            { key: '/init', desc: '创建 CLAUDE.md' },
            { key: '/memory', desc: '编辑 CLAUDE.md 文件' },
            { key: '/mcp', desc: '管理 MCP 服务器' },
            { key: '/hooks', desc: '管理钩子' },
            { key: '/skills', desc: '列出可用技能' },
            { key: '/agents', desc: '管理代理' },
            { key: '/chrome', desc: 'Chrome 集成' },
            { key: '/reload-plugins', desc: '热重载插件' }
          ]
        },
        {
          title: '特殊',
          rows: [
            { key: '/btw <question>', desc: '附带提问（无上下文）' },
            { key: '/plan [desc]', desc: '计划模式（+ 自动启动）' },
            { key: '/loop [interval]', desc: '调度周期性任务' },
            { key: '/voice', desc: '按键说话语音（20 种语言）' },
            { key: '/doctor', desc: '诊断安装问题' },
            { key: '/rc', desc: '启用远程控制' },
            { key: '/pr-comments [PR]', desc: '获取 GitHub PR 评论' },
            { key: '/stats', desc: '使用连续记录和偏好' },
            { key: '/insights', desc: '分析会话报告' },
            { key: '/desktop', desc: '在桌面应用中继续' },
            { key: '/remote-control', desc: '桥接终端到 claude.ai/code', addedAt: '2026-03-23' },
            { key: '/stickers', desc: '订购贴纸！' }
          ]
        }
      ]
    },
    {
      id: 'memory',
      title: '📁 记忆与文件',
      theme: 'memory',
      groups: [
        {
          title: 'CLAUDE.md 位置',
          rows: [
            { key: './CLAUDE.md', desc: '项目（团队共享）' },
            { key: '~/.claude/CLAUDE.md', desc: '个人（所有项目）' },
            { key: '/etc/claude-code/', desc: '管理（组织级）' }
          ]
        },
        {
          title: '规则与导入',
          rows: [
            { key: '.claude/rules/*.md', desc: '项目规则' },
            { key: '~/.claude/rules/*.md', desc: '用户规则' },
            { key: 'paths: frontmatter', desc: '路径特定规则' },
            { key: '@path/to/file', desc: '在 CLAUDE.md 中导入' }
          ]
        },
        {
          title: '自动记忆',
          rows: [
            { key: '~/.claude/projects/<proj>/memory/', desc: 'MEMORY.md + 主题文件，自动加载' }
          ]
        }
      ]
    }
  ],
  [
    {
      id: 'workflows',
      title: '🧠 工作流与技巧',
      theme: 'workflows',
      groups: [
        {
          title: '计划模式',
          rows: [
            { key: '[Shift][Tab]', desc: 'Normal → Auto → Plan' },
            { key: '--permission-mode plan', desc: '以计划模式启动' }
          ]
        },
        {
          title: '思考与努力度',
          rows: [
            { key: '[Alt][T]', desc: '切换思考开 / 关' },
            { key: '"ultrathink"', desc: '当前轮次最大努力度' },
            { key: '[Ctrl][O]', desc: '查看思考（详细模式）' },
            { key: '/effort', desc: '○ low · ◐ med · ● high', addedAt: '2026-03-23' }
          ]
        },
        {
          title: 'Git Worktrees',
          rows: [
            { key: '--worktree name', desc: '每个功能独立分支' },
            { key: 'isolation: worktree', desc: '代理在独立 worktree 中' },
            { key: 'sparsePaths', desc: '仅检出所需目录', addedAt: '2026-03-23' },
            { key: '/batch', desc: '自动创建 worktrees' }
          ]
        },
        {
          title: '语音模式',
          rows: [
            { key: '/voice', desc: '启用按键说话' },
            { key: '[Space] (hold)', desc: '录音，释放发送' },
            { key: '20 种语言', desc: 'EN, ES, FR, DE, CZ, PL…' }
          ]
        },
        {
          title: '上下文管理',
          rows: [
            { key: '/context', desc: '用量 + 优化建议' },
            { key: '/compact [focus]', desc: '带焦点压缩' },
            { key: 'Auto-compact', desc: '约 95% 容量时触发' },
            { key: '1M context', desc: 'Opus 4.6（Max / Team / Ent）' },
            { key: 'CLAUDE.md', desc: '压缩后保留！' }
          ]
        },
        {
          title: '会话高级操作',
          rows: [
            { key: 'claude -c', desc: '继续上次对话' },
            { key: 'claude -r "name"', desc: '按名称恢复' },
            { key: '/btw question', desc: '附带提问，无上下文开销' }
          ]
        },
        {
          title: 'SDK / 无头模式',
          rows: [
            { key: 'claude -p "query"', desc: '非交互式' },
            { key: '--output-format json', desc: '结构化输出' },
            { key: '--max-budget-usd 5', desc: '费用上限' },
            { key: 'cat file | claude -p', desc: '管道输入' }
          ]
        },
        {
          title: '调度与远程',
          rows: [
            { key: '/loop 5m msg', desc: '周期性任务' },
            { key: '/rc', desc: '远程控制' },
            { key: '--remote', desc: 'claude.ai 上的 Web 会话' }
          ]
        }
      ]
    },
    {
      id: 'config',
      title: '⚙️ 配置与环境',
      theme: 'config',
      groups: [
        {
          title: '配置文件',
          rows: [
            { key: '~/.claude/settings.json', desc: '用户设置' },
            { key: '.claude/settings.json', desc: '项目（共享）' },
            { key: '.claude/settings.local.json', desc: '仅本地' },
            { key: '~/.claude.json', desc: 'OAuth, MCP, 状态' },
            { key: '.mcp.json', desc: '项目 MCP 服务器' }
          ]
        },
        {
          title: '重要设置',
          rows: [
            { key: 'modelOverrides', desc: '模型选择器映射自定义 ID' },
            { key: 'autoMemoryDirectory', desc: '自定义记忆目录' },
            { key: 'worktree.sparsePaths', desc: '稀疏检出目录', addedAt: '2026-03-23' }
          ]
        },
        {
          title: '重要环境变量',
          rows: [
            { key: 'ANTHROPIC_API_KEY' },
            { key: 'ANTHROPIC_MODEL' },
            { key: 'CLAUDE_CODE_EFFORT_LEVEL', desc: 'low / med / high' },
            { key: 'MAX_THINKING_TOKENS', desc: '0 = 关闭' },
            { key: 'ANTHROPIC_CUSTOM_MODEL_OPTION', desc: '自定义 /model 条目' },
            { key: 'CLAUDE_CODE_PLUGIN_SEED_DIR', desc: '多个插件种子目录' }
          ]
        }
      ]
    }
  ],
  [
    {
      id: 'skills',
      title: '🔧 技能与代理',
      theme: 'skills',
      groups: [
        {
          title: '内置技能',
          rows: [
            { key: '/simplify', desc: '代码审查（3 个并行代理）' },
            { key: '/batch', desc: '大规模并行修改（5-30 worktrees）' },
            { key: '/debug [desc]', desc: '从调试日志排查问题' },
            { key: '/loop [interval]', desc: '周期性调度任务' },
            { key: '/claude-api', desc: '加载 API + SDK 参考' }
          ]
        },
        {
          title: '自定义技能位置',
          rows: [
            { key: '.claude/skills/<name>/', desc: '项目技能' },
            { key: '~/.claude/skills/<name>/', desc: '个人技能' }
          ]
        },
        {
          title: '技能前置元数据',
          rows: [
            { key: 'description', desc: '自动调用触发器' },
            { key: 'allowed-tools', desc: '跳过权限提示' },
            { key: 'model', desc: '覆盖技能模型' },
            { key: 'effort', desc: '覆盖努力度', addedAt: '2026-03-23' },
            { key: 'context: fork', desc: '在子代理中运行' },
            { key: '$ARGUMENTS', desc: '用户输入占位符' },
            { key: '${CLAUDE_SKILL_DIR}', desc: '技能自身目录' },
            { key: '!`cmd`', desc: '动态上下文注入' }
          ]
        },
        {
          title: '内置代理',
          rows: [
            { key: 'Explore', desc: '快速只读（Haiku）' },
            { key: 'Plan', desc: '计划模式研究' },
            { key: 'General', desc: '全工具，复杂任务' },
            { key: 'Bash', desc: '终端独立上下文' }
          ]
        },
        {
          title: '代理前置元数据',
          rows: [
            { key: 'permissionMode', desc: 'default / acceptEdits / dontAsk / plan' },
            { key: 'isolation: worktree', desc: '在 git worktree 中运行' },
            { key: 'memory: user|project', desc: '持久化记忆' },
            { key: 'background: true', desc: '后台任务' },
            { key: 'maxTurns', desc: '限制代理轮次' },
            { key: 'SendMessage', desc: '恢复代理（替代 resume）', addedAt: '2026-03-23' }
          ]
        }
      ]
    },
    {
      id: 'cli',
      title: '🖥️ CLI 与标志',
      theme: 'cli',
      groups: [
        {
          title: '核心命令',
          rows: [
            { key: 'claude', desc: '交互式' },
            { key: 'claude "q"', desc: '带提示' },
            { key: 'claude -p "q"', desc: '无头模式' },
            { key: 'claude -c', desc: '继续上次' },
            { key: 'claude -r "n"', desc: '恢复' },
            { key: 'claude update', desc: '更新' }
          ]
        },
        {
          title: '重要标志',
          rows: [
            { key: '--model', desc: '设置模型' },
            { key: '-w', desc: 'Git worktree' },
            { key: '-n / --name', desc: '会话名称' },
            { key: '--add-dir', desc: '添加目录' },
            { key: '--agent', desc: '使用代理' },
            { key: '--allowedTools', desc: '预批准' },
            { key: '--output-format', desc: 'json / stream' },
            { key: '--json-schema', desc: '结构化' },
            { key: '--max-turns', desc: '限制轮次' },
            { key: '--max-budget-usd', desc: '费用上限' },
            { key: '--console', desc: '通过 Anthropic Console 认证' },
            { key: '--verbose', desc: '详细模式' },
            { key: '--bare', desc: '最小化无头（无 hooks / LSP）', addedAt: '2026-03-23' },
            { key: '--channels', desc: '权限中继 / MCP 推送', addedAt: '2026-03-23' },
            { key: '--remote', desc: 'Web 会话' },
            { key: '--chrome', desc: 'Chrome' }
          ]
        }
      ]
    }
  ]
]

export const claudeCodeCheatSheetContent: CheatSheetContent = {
  meta: claudeCodeCheatSheetMeta,
  changelog: claudeCodeCheatSheetChangelog,
  footer: claudeCodeCheatSheetFooter,
  columns: claudeCodeCheatSheetColumns
}
