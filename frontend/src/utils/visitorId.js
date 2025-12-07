/**
 * 用户标识工具类
 * 用于生成和管理用户的唯一标识符
 */

/**
 * 生成UUID v4
 * @returns {string} UUID字符串
 */
export function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

/**
 * 获取用户标识符
 * @returns {string|null} 用户标识符
 */
export function getVisitorId() {
  // 首先尝试从localStorage获取
  let visitorId = localStorage.getItem('visitor_id');
  
  // 如果localStorage中没有，则尝试从cookie获取
  if (!visitorId) {
    visitorId = getCookie('visitor_id');
  }
  
  return visitorId;
}

/**
 * 设置用户标识符
 * @param {string} visitorId 用户标识符
 */
export function setVisitorId(visitorId) {
  // 存储到localStorage
  localStorage.setItem('visitor_id', visitorId);
  
  // 同时存储到cookie（一年有效期）
  setCookie('visitor_id', visitorId, 365);
}

/**
 * 确保用户标识符存在
 * @returns {string} 用户标识符
 */
export function ensureVisitorId() {
  let visitorId = getVisitorId();
  
  // 如果不存在，则生成新的
  if (!visitorId) {
    visitorId = generateUUID();
    setVisitorId(visitorId);
  }
  
  return visitorId;
}

/**
 * 获取Cookie值
 * @param {string} name Cookie名称
 * @returns {string|null} Cookie值
 */
function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return parts.pop().split(';').shift();
  }
  return null;
}

/**
 * 设置Cookie值
 * @param {string} name Cookie名称
 * @param {string} value Cookie值
 * @param {number} days 有效期（天）
 */
function setCookie(name, value, days) {
  const expires = new Date();
  expires.setTime(expires.getTime() + (days * 24 * 60 * 60 * 1000));
  document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/`;
}