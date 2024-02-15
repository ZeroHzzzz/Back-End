# Back-End
寒假软件设计T3 后端项目仓库

# 项目介绍
学生综合测评成绩管理系统

分为三大模块：Student，Counsellor, Square


# 基础实现
- 数据:
    - 计划采用redis作为缓存层读取数据，使用mongodb作为持久化数据储存库，定期同步数据，这样可以调高用户的查询速度，提高可用性
    - 使用rabbitmq作为用户消息队列，并使用websocket协议进行消息监听，同时解决当生产者大量产生数据时，消费者无法快速消费的问题，实现系统之间的双向解耦
- 第三方库的使用：
    - mongo-driver
        - go.mongodb.org/mongo-driver/mongo
    - Go-Redis
        - https://github.com/redis/go-redis
    - amqp
        - https://github.com/streadway/amqp
    - websocket
        - https://github.com/gorilla/websocket
    - gin
        - https://github.com/gin-gonic/gin
    - excelize
        - https://github.com/xuri/excelize/v2
    - ratelimit
        - https://github.com/juju/ratelimit
    - jwt
        - https://github.com/dgrijalva/jwt-go

# 项目规范
go语言项目
- 一般变量使用小驼峰camelCase命名方式，避免缩写，避免使用下划线“_”
- 包外可见的变量名称首字母大写
- 命名常量时候全部大写
- 结构体内变量采用大驼峰CamelCase命名方式，首字母大写，以便其他包可见
- 接口命名为动词+er的形式，例如：Reader、Writer、Formatter。这样的命名模式能够清晰地表示接口的功能

接口
- 统一使用小驼峰camelCase命名方式
- 使用明确的相应状态码，其标准如下：
    - 2**成功
        - 200 OK： 请求成功。
        - 201 Created： 请求已创建成功，通常用于 POST 请求创建资源。
        - 204 No Content： 请求成功，但响应中没有内容，用于更新或删除资源
    - 3** 重定向
        - 301 Moved Permanently： 资源被永久移动到新位置。
        - 302 Found： 资源被临时移动到新位置。
        - 304 Not Modified： 客户端缓存仍有效，可以直接使用
    - 4xx 客户端错误：
        - 400 Bad Request： 请求无效或参数错误。
        - 401 Unauthorized： 需要身份验证。
        - 403 Forbidden： 服务器理解请求，但拒绝执行。
        - 404 Not Found： 请求的资源不存在。
    - 5xx 服务器错误：
        - 500 Internal Server Error： 服务器遇到错误，无法完成请求。
        - 502 Bad Gateway： 充当网关或代理的服务器从上游服务器收到无效响应。
        - 503 Service Unavailable： 服务器当前无法处理请求，通常是临时性的。

- 错误码规定：
	- 101, "账号或密码错误"
	- 102, "没有足够权限"
	- 103, "未到接口开放时间"
	- 104, "上下文错误"
	- 201, "参数错误"
	- 203, "文件错误"
	- 204, "解码错误"
	- 301, "MongoDB Client 初始化失败"
	- 302, "MongoDB 数据库操作失败"
	- 401, "RabbitMQ 初始化失败"
	- 402, "RabbitMQ 声明队列失败"
	- 403, "RabbitMQ 声明交换机失败"
	- 404, "RabbitMQ 交换机绑定队列失败"
	- 405, "RabbitMQ 声明消费者失败"
	- 406, "RabbitMQ 生产信息失败"
	- 501, "Redis 初始化失败"
	- 601, "服务器内部错误"
	- 602, "连接异常"
	- 603, "当前访问人数过多，服务器限流"

- 辅导员信息样例

| UserId  | UserName | Grade | Profession |
|---------|----------|-------|------------|
| 202333 | 目猫老师     | 大一    | 机械         |

- 学生信息样例

| UserId | UserName | Grade | Profession | Class |
|--------|----------|-------|------------|-------|
| 202301 | 艾斯比      | 大一    | 计算机        | 2301  |

- 学生成绩表格样例

| userId | userName  | grade | profession | class | academicYear | 基本评定分D1 | 集体评定等级分（*） | 社会责任记实分（*） | 思政学习加减分 | 违纪违规扣分 | 学生荣誉称号加减分（*） | 智育平均学分绩点 | 体育课程成绩T1 | 体育竞赛获奖得分（*） | 早锻炼得分 | 文化艺术实践成绩M1（*） | 文化艺术竞赛获奖得分M2（*） | 寝室日常考核基本分 | “文明寝室”创建、寝室风采展等活动加分（*） | 寝室行为表现与卫生状况加减分 | 志愿服务分L2 | 实习实训L3 | 创新创业竞赛获奖得分（*） | 水平等级考试（*） | 社会实践活动C2（*） | 社会工作C3 |
|--------|-----------|-------|------------|-------|--------------|---------|------------|------------|---------|--------|--------------|----------|----------|-------------|-------|---------------|-----------------|-----------|------------------------|----------------|---------|--------|---------------|-----------|-------------|--------|
| 202222 | zerohzzzz | 大一    | 计算机        | 2306  | 2023-2024    | 12      | 123        | 123        | 123     | 123    | 123          | 123      | 123      | 123         | 34    | 34            | 34              | 34        | 34                     | 34             | 34      | 34     | 34            | 34        | 34          | 34     |


# 接口文档地址

https://apifox.com/apidoc/shared-61642332-f0fa-4a35-8853-99662e0d9bcf

# 路由
| 请求方法 | 路径                        | 描述                                 |
|----------|-----------------------------|--------------------------------------|
| GET      | /login/student              | 学生登录                             |
| GET      | /login/counsellor           | 辅导员登录                           |
| GET      | /api/ws/:UserID             | WebSocket 连接                       |
| PUT      | /api/student/profile/:UserID | 修改学生个人信息                     |
| POST     | /api/student/feedbackOradvice | 提交学生反馈或建议                   |
| GET      | /api/student/score          | 获取学生成绩                         |
| POST     | /api/student/submit/:UserID | 提交学生作业                         |
| GET      | /api/student/submit/status/:SubmissionID | 获取作业提交状态             |
| GET      | /api/student/submit/list   | 获取学生作业列表                     |
| POST     | /api/counsellor/:CounsellorID/cause | 添加事由                        |
| GET      | /api/counsellor/:CounsellorID/cause | 获取事由                          |
| POST     | /api/counsellor/access-time | 设置辅导员可预约的时间              |
| POST     | /api/counsellor/setannouncement | 设置公告                           |
| GET      | /api/counsellor/audit/list | 获取审核列表                         |
| PUT      | /api/counsellor/audit/review/single | 单个审核                        |
| PUT      | /api/counsellor/audit/review/bulk | 批量审核                          |
| GET      | /api/counsellor/audit/history | 获取审核历史                       |
| PUT      | /api/counsellor/information/correct/:UserID | 校正学生成绩                  |
| POST     | /api/counsellor/information/bulk-import/student | 批量导入学生信息              |
| POST     | /api/counsellor/information/bulk-import/counsellor | 批量导入辅导员信息          |
| POST     | /api/counsellor/information/bulk-import/mark | 批量导入成绩信息               |
| GET      | /api/counsellor/information/student | 获取学生信息                       |
| GET      | /api/square/annoucement     | 获取广场公告                         |
| POST     | /api/square/topic/new      | 发布新主题                           |
| PUT      | /api/square/topic           | 修改主题                             |
| GET      | /api/square/topic/list     | 获取主题列表                         |
| GET      | /api/square/topic          | 获取主题详情                         |
| POST     | /api/square/topic/replies  | 发布回复                             |
| GET      | /api/square/topic/replies  | 获取回复列表                         |
| GET      | /api/square/topic/views&likes | 获取浏览量和点赞量                 |
| PUT      | /api/square/topic/likes/reply | 点赞回复                             |
| PUT      | /api/square/topic/likes/topic | 点赞主题                             |
| DELETE   | /api/square/topic/delete/topic | 删除主题                           |
| DELETE   | /api/square/topic/delete/reply | 删除回复                           |


# 申报表部分
下列可申请项目结尾用（*）标识<br>
而且有分项的需要进行汇总计算得分
- 德育素质
    - 基本评定分D1
    - 记实加减分D2
        - 集体评定等级分（*）
        - 社会责任记实分（*）
        - 思政学习加减分
        - 违纪违规扣分
        - 学生荣誉称号加减分（*）

- 智育素质
    - 智育平均学分绩点

- 体育素质
    - 体育课程成绩T1
    - 课外体育活动成绩T2
        - 体育竞赛获奖得分（*）
        - 早锻炼得分

- 美育素质
    - 文化艺术实践成绩M1（*）
    - 文化艺术竞赛获奖得分M2（*）

- 劳育素质
    - 日常劳动分L1
        - 寝室日常考核基本分
        - “文明寝室”创建、寝室风采展等活动加分（*）
        - 寝室行为表现与卫生状况加减分
    - 志愿服务分L2
    - 实习实训L3

- 创新与实践素质
    - 创新创业成绩C1
        - 创新创业竞赛获奖得分（*）
        - 水平等级考试（*）
    - 社会实践活动C2（*）
    - 社会工作C3

通过观察文档我们可以看出，每个模块只有一个父项有子项，那么我们就可以通过在父项处加入F-开头的字段进行标识，在子相处加入L-开头字段命名表示