需要先部署的服务：Redis,MySQL,Minio



MinIO社区版在2025-05-24T17-08-30Z这个版本之后，移除了控制台大部分管理功能，想要完整minio，请安装旧版本。

先在usr/local路径下创建minio-data文件夹以及config文件夹

-e "MINIO_ROOT_USER=xxx" \    设置minio的用户名
-e "MINIO_ROOT_PASSWORD=xxxxxxxx" \    设置minio的用户密码

```
mkdir -p /opt/minio/data
mkdir -p /opt/minio/config

docker run -d \
  --name minio \
  -p 9000:9000 \
  -p 9001:9001 \
  -e "MINIO_ROOT_USER=minio" \
  -e "MINIO_ROOT_PASSWORD=minio123" \
  -v /usr/local/minio-data:/data \
  -v /usr/local/minio-config:/root/.minio \
  --restart unless-stopped \
  minio/minio:RELEASE.2025-04-22T22-12-26Z server /data --console-address ":9001"
```





nginx的conf配置文件

```
worker_processes  1;

events {
    worker_connections  1024;
}


http {
    include       mime.types;
    default_type  application/octet-stream;
    client_max_body_size 30M;

    # 禁用 keepalive 超时（避免连接过早关闭）
    keepalive_timeout 0;

    sendfile        on;

    server {
        listen       80;
        server_name  localhost;

        location / {
            root   dist;
            index  index.html index.htm;
	    try_files $uri $uri/ /index.html;
        }
	
	location /wallpaper {
	     #proxy_pass http://59.153.164.121:8088;
		proxy_pass http://localhost:8888;
	     # ...其他配置
             #proxy_set_header Host $http_host; # 超级重要
	}

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    }
}
```

