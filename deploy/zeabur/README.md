# Dujiao-Next on Zeabur

This guide contains two deployment modes.

## Mode A: one-service fullstack deployment

The root `Dockerfile` builds the configured user/admin frontends and embeds them into the Go API binary with the `fullstack` build tag. By default it uses:

- https://github.com/wubowubohh-jpg/user
- https://github.com/wubowubohh-jpg/admin

This is the simplest Zeabur deployment:

- PostgreSQL
- Redis
- One Dujiao-Next service

The service exposes:

- user storefront at `/`
- admin frontend at `web.admin_path`, default `/admin`
- API at `/api/v1`
- uploads at `/uploads`

### 1. Domains

Prepare one public domain:

```text
__SHOP_DOMAIN__   # example: shop.example.com
```

Bind it directly to the Dujiao-Next API/fullstack service on port `8080`.

If you want the admin path to be harder to guess, set it in `/app/config.yml`:

```yaml
web:
  admin_path: "/console-private"
```

Then the admin frontend is available at:

```text
https://__SHOP_DOMAIN__/console-private
```

### 2. Database and Redis

Create PostgreSQL and Redis from Zeabur templates. Copy their private hostnames from Zeabur Networking.

### 3. Dujiao-Next service

Deploy this GitHub repository directly. Zeabur should build it with the root `Dockerfile`.

Service port:

```text
8080
```

Mount volumes:

```text
/app/uploads
/app/logs
```

Use Zeabur Config Editor to mount `api.config.yml.example` as:

```text
/app/config.yml
```

For Mode A, replace both `__SHOP_DOMAIN__` and `__ADMIN_DOMAIN__` in CORS with the same public domain, unless you later place admin on a separate host.

Recommended environment variables:

```env
TZ=Asia/Shanghai
DJ_DEFAULT_ADMIN_USERNAME=admin
DJ_DEFAULT_ADMIN_PASSWORD=replace-with-a-strong-first-login-password
```

### 4. Checks

```text
https://__SHOP_DOMAIN__
https://__SHOP_DOMAIN__/api/v1/public/config
https://__SHOP_DOMAIN__/sitemap.xml
https://__SHOP_DOMAIN__/robots.txt
https://__SHOP_DOMAIN__/admin
https://__SHOP_DOMAIN__/health
```

## Mode B: split frontend deployment

Use this mode only when you want the official frontend Docker images as separate services or need independent scaling. It deploys:

- PostgreSQL
- Redis
- API
- User frontend
- Admin frontend
- Nginx edge proxy

### 1. Domains

Prepare two domains:

```text
__SHOP_DOMAIN__   # example: shop.example.com
__ADMIN_DOMAIN__  # example: admin.example.com
```

Bind both domains to the Nginx service, not directly to API/User/Admin.

### 2. Database and Redis

Create these services from Zeabur templates:

- PostgreSQL
- Redis

Open each service's `Networking` tab and copy the `Private` hostname. Zeabur documents this private hostname flow in its Private Networking docs.

### 3. API service

Use one of these two approaches:

- Deploy this GitHub repository directly. Zeabur should build it with the root `Dockerfile`.
- Or deploy the official Docker image, for example `dujiaonext/api:latest`.

API port:

```text
8080
```

Mount volumes:

```text
/app/uploads
/app/logs
```

If you switch back to SQLite, also mount:

```text
/app/db
```

Use Zeabur Config Editor to mount `api.config.yml.example` as:

```text
/app/config.yml
```

Replace all placeholders before starting the service.

Recommended API environment variables:

```env
TZ=Asia/Shanghai
DJ_DEFAULT_ADMIN_USERNAME=admin
DJ_DEFAULT_ADMIN_PASSWORD=replace-with-a-strong-first-login-password
```

After the first successful admin login, change the admin password and remove or rotate the bootstrap password.

### 4. User frontend service

Create a Docker Image service:

```text
dujiaonext/user:latest
```

Port:

```text
80
```

Environment:

```env
TZ=Asia/Shanghai
```

Do not bind the public shop domain directly to this service.

### 5. Admin frontend service

Create a Docker Image service:

```text
dujiaonext/admin:latest
```

Port:

```text
80
```

Environment:

```env
TZ=Asia/Shanghai
```

Do not bind the public admin domain directly to this service.

### 6. Nginx service

Create a Docker Image service:

```text
nginx:1.27-alpine
```

Port:

```text
80
```

Use Zeabur Config Editor to mount `nginx.conf.example` as:

```text
/etc/nginx/conf.d/default.conf
```

Replace the placeholders:

```text
__SHOP_DOMAIN__
__ADMIN_DOMAIN__
__API_PRIVATE_HOST__
__USER_PRIVATE_HOST__
__ADMIN_PRIVATE_HOST__
```

Then bind both public domains to this Nginx service.

### 7. Checks

After deployment, check:

```text
https://__SHOP_DOMAIN__
https://__SHOP_DOMAIN__/api/v1/public/config
https://__SHOP_DOMAIN__/sitemap.xml
https://__SHOP_DOMAIN__/robots.txt
https://__ADMIN_DOMAIN__
https://__ADMIN_DOMAIN__/api/v1/public/config
```

Expected API health check:

```text
https://__SHOP_DOMAIN__/health
```

If `/health` is needed from the public domain, add an Nginx location for it. The public storefront usually does not need this route.

## Payment callbacks

Use the public shop domain for payment callback URLs, for example:

```text
https://__SHOP_DOMAIN__/api/v1/payments/callback
https://__SHOP_DOMAIN__/api/v1/payments/webhook/dujiaopay
https://__SHOP_DOMAIN__/api/v1/payments/webhook/paypal
https://__SHOP_DOMAIN__/api/v1/payments/webhook/stripe
```

## References

- Dujiao-Next Docker Compose deployment: https://dujiao-next.com/deploy/docker-compose
- Dujiao-Next config: https://dujiao-next.com/config/config-yml
- Zeabur Private Networking: https://zeabur.com/docs/en-US/deploy/networking/private-networking
- Zeabur Config Editor: https://zeabur.com/docs/en-US/operations/data/config-file-management
- Zeabur Volumes: https://zeabur.com/docs/en-US/data-management/volumes
