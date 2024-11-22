# Clipboard-Go

Share clipboard between 💻 Linux and 📱 iOS

## Install

```bash
paru -Syu clipboard-go
systemctl --user daemon-reload
systemctl --user enable --now clipboard-go.service
```

### iOS Shortcuts
- **Copy**: [Copy Clipboard](https://www.icloud.com/shortcuts/82448695a1b8407a90e6abceee89ffac)
- **Paste**: [Paste Clipboard](https://www.icloud.com/shortcuts/da413599218348ea94a2ad9f1c0fa0ab)

## API

### Overview
The API provides functions to interact with the clipboard, including retrieving and setting clipboard contents. Requests are sent over HTTP in JSON format.

### Endpoints

#### 1. Get Clipboard Content
- **Request Method**: `GET`
- **Request Path**: `/`

**Request Headers**:
| Parameter   | Type   | Description              |
|-------------|--------|--------------------------|
| Timestamp   | string | UNIX timestamp (seconds) |
| MD5         | string | Authentication MD5 hash   |

**Response**:
- **Success**: 200 OK (returns clipboard content)
- **Error**:
  - 403 Forbidden
  - 500 Internal Server Error

#### 2. Set Clipboard Content
- **Request Method**: `POST`
- **Request Path**: `/`

**Request Headers**:
| Parameter   | Type   | Description              |
|-------------|--------|--------------------------|
| Timestamp   | string | UNIX timestamp (seconds) |
| MD5         | string | Authentication MD5 hash   |

**Request Body**:
- **Content Type**: `application/json`
- **Example**:
  ```json
  {
      "data": "New clipboard content"
  }
  ```

**Response**:
- **Success**: 204 No Content
- **Error**:
  - 400 Bad Request
  - 403 Forbidden
  - 500 Internal Server Error

### Authentication
- The client must include `Timestamp` and `MD5` in the request headers.
- The server validates the request's authenticity using the MD5 hash, calculated as follows:
  ```text
  MD5(timestamp + authPassword)
  ```

## 🌟Stargazers over time
[![Stargazers over time](https://starchart.cc/Xarth-Mai/Clipboard-Go.svg?variant=adaptive)](https://starchart.cc/Xarth-Mai/Clipboard-Go)
