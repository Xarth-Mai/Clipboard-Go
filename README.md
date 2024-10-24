# Clipboard-Go

在 💻Linux 与 📱iOS 之间共享剪贴板

## Install
```
paru -Syu clipboard-go
systemctl --user daemon-reload
systemctl --user start clipboard-go.service
systemctl --user enable clipboard-go.service
```

### iOS 快捷指令
- **复制**: [复制剪贴板](https://www.icloud.com/shortcuts/82448695a1b8407a90e6abceee89ffac)
- **粘贴**: [粘贴剪贴板](https://www.icloud.com/shortcuts/da413599218348ea94a2ad9f1c0fa0ab)

## API 

### 概述
该 API 提供了与剪贴板交互的功能，包括获取和设置剪贴板内容。请求通过 HTTP 协议发送，采用 JSON 格式。

### 端点

#### 1. 获取剪贴板内容
- **请求方法**: `GET`
- **请求路径**: `/`
  
**请求头**:
| 参数       | 类型   | 描述                    |
|------------|--------|-------------------------|
| Timestamp  | string | UNIX 时间戳 (秒)       |
| MD5        | string | 对应的 MD5 哈希值      |

**响应**:
- **成功**: 200 OK (返回剪贴板内容)
- **错误**: 
  - 403 Forbidden
  - 500 Internal Server Error

#### 2. 设置剪贴板内容
- **请求方法**: `POST`
- **请求路径**: `/`

**请求头**:
| 参数       | 类型   | 描述                    |
|------------|--------|-------------------------|
| Timestamp  | string | UNIX 时间戳 (秒)       |
| MD5        | string | 对应的 MD5 哈希值      |

**请求体**:
- **内容类型**: `application/json`
- **示例**:
  ```json
  {
      "data": "新的剪贴板内容"
  }
  ```

**响应**:
- **成功**: 204 No Content
- **错误**:
  - 400 Bad Request
  - 403 Forbidden
  - 500 Internal Server Error

### 身份验证
- 客户端需在请求头中包含 `Timestamp` 和 `MD5`。
- 服务器通过 MD5 哈希验证请求有效性，哈希计算如下：
```
MD5(timestamp + authPassword)
```

### 状态码
- **200 OK**: 请求成功
- **204 No Content**: 请求成功但无返回内容
- **400 Bad Request**: 请求格式不正确
- **403 Forbidden**: 身份验证失败
- **500 Internal Server Error**: 服务器内部错误

## 示例请求

### 获取剪贴板内容
```http
GET / HTTP/1.1
Host: localhost:8777
Timestamp: 1635434567
MD5: 1a79a4d60de6718e8e5b326e338ae533
```

### 设置剪贴板内容
```http
POST / HTTP/1.1
Host: localhost:8777
Timestamp: 1635434567
MD5: 1a79a4d60de6718e8e5b326e338ae533
Content-Type: application/json

{
    "data": "新的剪贴板内容"
}
```
