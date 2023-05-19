# gohtran
Golang 版本的 Htran 工具

## 能做什么
实现不同端口串联在一起

## 例如：
你希望通过A机器的2222端口链接至机器B的3389端口上，假设b的地址为192.168.1.100，那么你可以这样做

`gohtran -tran 2222 192.168.1.100:3389`

通过win+R唤出mstsc,输入本机地址127.0.0.1:2222便可以连接到b机器上进行远程控制电脑

linux同理

`ssh -p 2222 username@127.0.0.1`

便可以实现通过链接本地的2222端口，连接到192.168.1.100机器的3389上

## 安装
```bat
git clone git@github.com:7134g/gohtran.git
go build .
```

## 主要功能介绍
1. 本地监听转发

   将本机的2222端口与本机的3333端口绑定，向2222端口写入的数据可以从3333中读取,反之亦然
   - `gohtran -listen 2222 3333`
2. 转发到远程主机

   实现链接127.0.0.1:2222时候，实际访问的是192.168.1.100:3306端口
   - `gohtran -tran 2222 192.168.1.100:3306`
3. 反向连接转发

   同时链接A（192.168.1.101）和B（192.168.1.100）两台机器的端口，
   实现其他机器远程链接A机器2222端口时候，实际上链接的是B机器的3389端口
   - `gohtran -slave 192.168.1.101:2222 192.168.1.100:3389`


根据上述功能就可以实现无限串联所有可到达的机器
```text
+-----------------------------help information--------------------------------+
usage: "-listen port1 port2" #example: "gohtran -listen 8888 3389"
       "-tran port1 ip:port2" #example: "gohtran -tran 8888 1.1.1.1:3389"
       "-slave ip1:port1 ip2:port2" #example: "gohtran -slave 127.0.0.1:3389 1.1.1.1:8888"
       "-e enable gzip and aes functionality
       "-aes enable aes functionality, parameters is key, defaults to 16 bits
       "-gzip enable gzip functionality
       "-h -help program documentation
       "-s -silent silent mode,no information is displayed
       "-log output transferred data to file
============================================================
If you see start transmit, that means the data channel is established
```

#### 注意事项
若开启了aes加密或者gzip压缩，则至少需要运行两个程序才可以正常加解密

例如：
   1. A（192.168.1.101）机器开启listen模式`gohtran -listen 2222 3333 -aes`
   2. B（192.168.1.100）机器开启slave模式`gohtran -slave 192.168.1.101:3333 192.168.1.100:3389 -aes`
   3. 这时候访问A机器端口2222才可以正常解析数据
   
