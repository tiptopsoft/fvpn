#/bin/sh
# 查询进程
process=`ps -ef | grep "node" | grep -v grep | awk '{print $2}'`

# 如果进程存在，则结束该进程
if [[ -n "$process" ]]; then
  kill $process
fi



#启动node
chmod +x fvpn

nohup ./fvpn node &

#join network
./fvpn join 96141f705c81ccc1