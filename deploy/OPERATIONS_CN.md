# ukapi.cc 运维说明

本文档对应 2026-08-04 已上线的 `sub2api` Docker 部署。

## 访问入口

- 主站地址: `https://ukapi.cc`
- 健康检查: `https://ukapi.cc/health`
- 登录页: `https://ukapi.cc/login`
- 管理后台: `https://ukapi.cc/admin/dashboard`
- `www.ukapi.cc` 已配置为 301 跳转到 `https://ukapi.cc`

## 当前部署结构

- Cloudflare: `Full (strict)` + 橙云代理
- Nginx: 监听 `80/443`
- Sub2API: 容器内服务监听 `127.0.0.1:8080`
- PostgreSQL: Docker 容器
- Redis: Docker 容器

## 服务器关键信息

- 服务器 IP: `154.19.229.193`
- 项目目录: `/root/sub2api`
- 部署目录: `/root/sub2api/deploy`
- Nginx 证书目录: `/etc/nginx/ssl/sub2api`
- Nginx 站点配置: `/etc/nginx/sites-available/ukapi.cc.conf`

## 管理员账号

- Email: `admin@ukapi.cc`
- Password: `e023873215c6202802203257697923e2`

首次登录建议立刻修改管理员密码，并检查后台基础设置。

## 邮件配置

- Provider: Resend SMTP
- SMTP host: `smtp.resend.com`
- SMTP port: `587`
- SMTP username: `resend`
- TLS: `STARTTLS`
- From: `openapi@ukapi.cc`
- Frontend URL: `https://ukapi.cc`

已完成测试:

- SMTP connection: successful
- Test recipient: `delivered@resend.dev`
- Send test email: successful

Resend API key 未写入仓库文件，只保存在应用数据库的加密 SMTP 密码字段中。

## 关键配置现状

- 对外开放端口: `80`, `443`
- `8080` 已收回为本机回环，不再暴露公网
- Cloudflare Origin Certificate 已部署到服务器
- Nginx 反向代理目标: `http://127.0.0.1:8080`

## 常用检查命令

```bash
docker ps
ss -lntp | grep -E ':(80|443|8080)\b'
curl -fsS http://127.0.0.1:8080/health
curl -fsS https://ukapi.cc/health
nginx -t
systemctl status nginx --no-pager
```

## 常用重载命令

```bash
systemctl restart nginx
cd /root/sub2api/deploy
SUB2API_IMAGE=sub2api-fork SUB2API_IMAGE_TAG=git-1c67cc5e1f5b docker compose -f docker-compose.local.yml -f docker-compose.fork.yml up -d --remove-orphans --no-build
```

## 备份命令

```bash
cd /root/sub2api/deploy
COMPOSE_PROJECT_NAME=deploy BACKUP_ROOT=/root/sub2api/deploy/backups bash backup.sh
```

备份目录:

- `/root/sub2api/deploy/backups`

## 仓库内相关文件

- `deploy/nginx/ukapi.cc.conf`: 当前线上 Nginx 配置模板
- `tools/configure_cloudflare_nginx.py`: 证书上传和 Nginx 配置脚本
- `tools/harden_sub2api_bind.py`: 将 `8080` 收回到 `127.0.0.1` 的脚本

## 注意事项

- 这台服务器内存较小，直接重建前端镜像容易 OOM。
- 重新拉起服务时，优先复用现成镜像并带上 `--no-build`。
- 如果后续更换域名或证书，需要同步修改 Nginx 配置和 Cloudflare 源站证书。
