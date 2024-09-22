本仓库实现自建飞书机器人与本地 [ragflow](https://github.com/infiniflow/ragflow) 服务交互，实现本地知识库对话，支持本地文档预览。

功能实现如下：
- 配置飞书机器人 webhook。
- 配置 ragflow 服务。
- @机器的消息实现自动饮用，并且搭配 ragflow 进行知识库对话。
- ragflow 返回的文档支持预览，并且在消息最后分条显示。
- 回复的消息支持 markdown 格式。