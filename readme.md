# Docker commander
basically I was design this to handle gitlab webhook to run command.  
Why? because pipeline have to pull image and run so why do I have to wait for it?  

## Example config
```dotenv
# need colon prefix due it's host:port but host is optional
LISTEN=:3000
IMAGE_nginx=docker compose restart nginx
IMAGE_db=docker compose restart database && curl http://notify/something
TOKEN=1234
```
if we send 
```http request
POST /nginx
X-Gitlab-Token 1234
```
or
```http request
GET /nginx?token=1234
```
it will run `docker compose restart nginx`