# Buffet: Ifttt Service for collect any link  : -)

Collect anything you like to your own database. ( eg: airtable & spreadsheet & ... ).


Buffet, A IFTTT service。 

`IF any new item from pocket, then extract meta information and save to airtable`.




# 构架 Architecture

A Simple HTTP Server + Spider

Http Server 用来接受 Ifttt 发送过来的 Action， 然后启动一个爬虫，解析好网页信息之后 推送到 Ifttt 。

Task 使用循环列表（要什么数据库，直接保存在内存!）保存最近的 100 个 任务。





# 帮帮可怜的乌干达儿童！

Contribute Steps:

-  Make ensure you have golang 1.10 & dep installed 
-  `mkdir -p $GOPATH/src/github.com/yowenter/buffet`
-  `git clone $(your forked git) $GOPATH/src/github.com/yowenter/buffet`
-  checkout a branch and write code !

Ifttt Develop Docs: `https://platform.ifttt.com/docs`


# Usage

- go run main.go
- open the `postman_buffet.json` with postman





# RoadMap

- [ ] Ifttt Trigger
- [X] 支持 Douban
- [ ] 支持 Github
- [ ] 支持 DianPing 




