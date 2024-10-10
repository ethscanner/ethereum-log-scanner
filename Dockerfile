FROM golang AS base

WORKDIR /app

COPY . .

# 构建两个main程序
RUN go build -o main

RUN chmod +x ./main

# 安装supervisor
RUN apt update && apt install -y supervisor

# 复制supervisor配置文件
COPY conf /etc/supervisor/conf.d

# 使用supervisor作为容器的入口点
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/supervisord.conf"]