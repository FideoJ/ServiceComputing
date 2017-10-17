## Introdution
课程《服务计算》作业：golang实现selpg。
相关链接：[开发 Linux 命令行实用程序](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)
## Install
```
go get github.com/FideoJ/ServiceComputing/hw1
```

## Help
```
$GOPATH/bin/hw1 -h
```

## Notes
使用`cat -n`代替`lp -dlp0`等，详见注释。

## Examples
```
$ $GOPATH/bin/hw1 -s2 -e4 -l2 /etc/passwd
bin:x:2:2:bin:/bin:/usr/sbin/nologin
sys:x:3:3:sys:/dev:/usr/sbin/nologin
sync:x:4:65534:sync:/bin:/bin/sync
games:x:5:60:games:/usr/games:/usr/sbin/nologin
man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
hw1: done


$ $GOPATH/bin/hw1 -s2 -e4 -l2 -dlp0 /etc/passwd
     1	bin:x:2:2:bin:/bin:/usr/sbin/nologin
     2	sys:x:3:3:sys:/dev:/usr/sbin/nologin
     3	sync:x:4:65534:sync:/bin:/bin/sync
     4	games:x:5:60:games:/usr/games:/usr/sbin/nologin
     5	man:x:6:12:man:/var/cache/man:/usr/sbin/nologin
     6	lp:x:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
hw1: done


$ echo "haha" | hw1 -s1 -e1 -l2 -dlp0
     1	haha
hw1: done
```