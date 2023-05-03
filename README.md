# Pandora-GO

这是一个基于[潘多拉 (Pandora)](https://github.com/pengzhile/pandora)的网页API的SDK。

这个项目不是为了复刻潘多拉的功能，而是为了能在golang程序中实现网页中那种对话的效果，方便开发各种应用。

## 示例程序中用到的环境变量配置

**CHAT_TOKEN**: 通过这个网站获取https://chat-api.zhile.io/auth
**CHAT_URL**: 目前这个可用https://chat-api.zhile.io/api

## 目前支持的功能

1. [x] 获取模型列表
2. [x] 开启新的对话
3. [x] 提交一个会话消息并等待返回结果
4. [x] 重命名一个会话
5. [x] 删除一个会话
6. [ ] 重新生成消息
