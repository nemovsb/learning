!#/bin/bash
docker pull redis
docker run --name redis-test-instance -p 6379:6379 -d redis
sudo apt install redis-tools
redis-cli -h localhost -p 6379