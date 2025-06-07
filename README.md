# LanDropUi

## 项目简介
<img width="1024" alt="17311749308526_ pic" src="https://github.com/user-attachments/assets/b493c286-1be1-4a39-ae7e-c5f8ba7340db" />
<img width="1024" alt="17321749308534_ pic" src="https://github.com/user-attachments/assets/d09e06d0-1ad0-4786-9c75-5ec18f3e823c" />



LanDropUi 是一个基于 Wails 框架开发的桌面应用程序，旨在提供局域网内的文件传输功能。它包含两个独立的应用程序：`receive`（接收端）和 `send`（发送端），允许用户方便快捷地在同一局域网内传输文件。

## 功能特性

- **文件发送 (`send`)**：用户可以选择一个或多个文件，并通过 TCP 连接发送到局域网内的接收端。
- **文件接收 (`receive`)**：监听指定端口，接收来自发送端的文件，并显示传输进度。
- **实时进度显示**：在文件传输过程中，实时更新传输进度。
- **跨平台**：基于 Wails，支持 Windows、macOS平台。
- **支持多文件同时发送/接收**：可手动选择传输/接收的文件数量。

## 技术栈

- **后端**：Go
- **前端**：Vue.js
- **桌面应用框架**：Wails v2

## 使用方式
在Releases中下载对应系统的客户端即可


