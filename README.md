# 深度分页Demo

使用Redis zset解决深度分页查询的性能问题

运行步骤
```
cd docker-compose
docker-compose up -d

cd ..
make 

# 生成5000w测试数据
./posts-list gen --batch-size=10000 --count=50000000 --cocurrent=4

# 构建zset索引
./posts-list build-index -c 4

# 启动服务器
./posts-list serve

# 查询
curl localhost:8080/api/posts?size=100&page=10000
```