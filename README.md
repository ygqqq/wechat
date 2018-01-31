# 项目简介

这个小项目是采用纯前后端分离的模式来做的一个模仿微信的webapp，前台基于vuejs 2.x实现，后台使用go语言提供restful api。为了实现即时通信，前段采用了html5的websocket。

由于是前后端分离，所以前端项目和后端项目(位于server目录下)完全可以分别部署。前端的部署比较简单，直接`npm install`安装所需依赖就可以，后端的部署比较复杂，后端由于是基于go语言实现，所以要配置go环境，另外项目中还用到了mongodb作为主存储，采用redis作为服务端数据缓存，采用kafka来作为消息中间件，所以部署起来有些麻烦。

基于这些原因，为了方便前后端协同开发，后端采用docker进行集成，所有的后端环境(包括go、redis、mongodb、kafka、zookeeper等)全部运行在一个docker容器中，并暴露相关api接口和端口给前端调用。

由于这个小项目只是一时兴起，也只是用来练手，再加上年底这段时间工作较忙，只能晚上抽时间来写一点，所以有非常多粗糙和不完善的地方，并且后端的架构也没有做成分布式的(基于这个项目再改造成分布式的也不麻烦)。为了方便开发和快速迁移，所有的后端环境全部集成在一个docker容器中了，并且所有的结点(包括redis、mongodb、kafka和zookeeper)都是单点，并且是同一个容器中。

这样做也有个好处就是非常方便部署，哪怕换了环境和电脑，只要安装了docker，把代码从git上拉下来，再根据写好的dockerfile构建镜像，再把用于调用api接口url的配置文件写好，就可以几条指令就将整个项目跑起来，整个部署环节轻松无痛，可以迅速进入开发状态，非常方便协同开发。

为了方便部署和开发，建议使用mac或者ubuntu系统，这样可以直接在本机既部署前端，又部署后端。如果想把代码拉下来在本地跑一下看下效果的话，可以依照下节所讲的步骤来。不过这里说的是开发环境的部署，如果想部署到服务器上，则需要对前端和后端分别打包编译和部署。

# 项目部署

首先要安装docker，具体安装过程可以参考docker官网，针对不同的系统进行安装，这里不再赘述。

然后是拉取代码，假设是在`/home/ygq`目录下
```bash
git clone https://github.com/ygqqq/wechat.git
```
然后安装前端所需的一些库，前提是已经配置好node环境，这里就不再讲如何配置node环境了
```bash
cd wechat
npm install     #如果网速不给力的话，建议使用淘宝的npm镜像源
```

然后需要构建docker镜像，这一步根据网速不同，需要的时间也不同，不过一般会比较慢，10-30分钟都有可能，因为要安装很多东西来配置环境
``` bash
cd server
sudo docker build -t mygo .     # 一定得在server目录里面构建镜像，因为前面的指令.指定了在当前目录找dockerfile
```
经过漫长的等待，如果一切顺利的话，镜像应该是构建成功了，下面就需要运行docker容器
```bash
# 这里是将docker运行在后台，并且将本地的server目录和docker容器中的/app目录建立映射关系，这样可以很方便开发，比如后台代码如果修改了，只需要重启go程序即可。-v后面的路径一定要正确　-v 你的server路径:/app   -p 你本地端口:docker容器端口　　后面的mygo是刚才构建的镜像名字，可以自定义

sudo docker run -d -p 8000:8000 -v /home/ygq/wechat/server/:/app --name=mygo mygo
```
然后docker容器就在后台运行了起来，如果想进入容器里面进行一些操作，比如启动go程序、查看下mongodb数据库、redis运行状态等等的
```bash
# 这里的mygo是刚才 --name=mygo指定的容器名称，可以自定义
sudo docker exec -it mygo /bin/bash
```
进入容器之后，redis、mongodb、kafka、zookeeper、go程序都已经默认自启动了，可以使用`ss -tnl`看到各自监听在自己默认的端口上(go是监听在8000端口)，如果修改了后台代码，可以自行重启go程序。

完成了上述操作之后，后端就部署差不多了，虽然有点麻烦，但好在只需要部署一次就够，以后开发就很方便了。剩下的就是修改下配置文件：
```javascript
// 为了协同开发，所以将一些api接口的调用url写到了配置文件中，需要自行配置
// 请参考wechat/config/local.config.example.js文件，拷贝一份，重命名为local.config.js,内容如下：
module.exports = {
    //本地后台api调用的地址和端口
    apiUrl : "http://127.0.0.1:8000",  
    //websocket调用接口     
    wsUrl : "ws://127.0.0.1:8000/ws"
}
//如果docker容器和你前端是跑在同一台机器上，就默认配置就可以，否则需要将ip改为docker容器所在的ip，注意不是docker容器自身的ip，而是宿主机的ip
```
修改完配置文件之后，一切完成，只差最后一步了
```bash
cd /home/ygq/wechat
npm run dev
```


