# Simp

1. 我喜欢TAF Web化的部署方式。
2. 我喜欢TRPC 简洁的开发模式
3. TAF太重了，作为使用TAFNode进行日常开发的程序员，也会经常遇到自己不能hold住的问题。

我想让所有事情都简单化。能不能将所有东西都弄简单一点。不想背八股文。

1. 我还是很喜欢TAF的Web化部署方式，重点是，我爱JavaScript！
2. 部署平台不能过多的涉及业务。TAF平台里所展示的一些东西其实和业务也是有绑定的，这使得开发者得花过多的时间去研究和了解。
3. 通信我觉得都采用HTTP就可以了，流量带宽的，很多时候并不需要去考虑。

配置采用 yaml的形式

````yaml
server:
  name: SimpServer
  port: 8080
  type: main
  staticPath: static # 静态资源
  storage:  mysql@3306...... # 存储
  proxy:
    - server:
        type: fass # 选择Fass时默认代表本地，并且！，冷启动！，不能包含定时器等玩意，
        name: FassServer
    - server:
        type: http
        port: 9091
        name: StudentServer
````