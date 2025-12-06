"""认证服务"""
from utils.jwt_utils import verify_password, create_access_token


def authenticate_admin(password: str) -> dict:
    """
    验证管理员密码并返回token
    :param password: 密码
    :return: {"success": bool, "token": str, "message": str}
    """
    if verify_password(password):
        token = create_access_token(data={"role": "admin"})
        return {
            "success": True,
            "token": token,
            "message": "登录成功"
        }
    else:
        return {
            "success": False,
            "token": None,
            "message": "密码错误"
        }
