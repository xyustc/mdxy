"""User-Agent 解析工具"""
from typing import Dict
from user_agents import parse


def parse_user_agent(user_agent_string: str) -> Dict[str, str]:
    """
    解析 User-Agent 字符串
    返回: {device_type, os, browser}
    """
    if not user_agent_string:
        return {
            "device_type": "Unknown",
            "os": "Unknown",
            "browser": "Unknown"
        }
    
    try:
        user_agent = parse(user_agent_string)
        
        # 设备类型
        if user_agent.is_mobile:
            device_type = "Mobile"
        elif user_agent.is_tablet:
            device_type = "Tablet"
        elif user_agent.is_pc:
            device_type = "PC"
        else:
            device_type = "Other"
        
        # 操作系统
        os_family = user_agent.os.family if user_agent.os.family else "Unknown"
        os_version = user_agent.os.version_string if user_agent.os.version_string else ""
        os = f"{os_family} {os_version}".strip()
        
        # 浏览器
        browser_family = user_agent.browser.family if user_agent.browser.family else "Unknown"
        browser_version = user_agent.browser.version_string if user_agent.browser.version_string else ""
        browser = f"{browser_family} {browser_version}".strip()
        
        return {
            "device_type": device_type,
            "os": os,
            "browser": browser
        }
    except Exception:
        return {
            "device_type": "Unknown",
            "os": "Unknown",
            "browser": "Unknown"
        }
