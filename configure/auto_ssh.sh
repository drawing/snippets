# usage auto_login.sh configure
# configure file:
#	servername	username	passwd	127.0.0.1	command
#	push-1-1	testname	123456	127.0.0.1	command
#
# add to system:
# ~/.bashrc
#	. auto_ssh.sh init configure
#
# usage:
# 	auto servername
#

script_name=auto_ssh.sh

[ $# -le 1 ] && { echo "Usage: $script_name command configure"; return;}

target=$1
config=$2

if [ "$target" == "init" ];
then
	alias auto="$script_name login $2";
	alias aucp="$script_name scp $2";

	name_list=`awk '{printf("%s ", $1)}' $config`;
	complete -W "$name_list" auto
	return
fi


if [ "$target" == "login" ];
then
	host=$3

	read name user pass ip cmd <<< $(awk -v key="$host" '{if ($1==key)print $1,$2,$3,$4,$5};' $config);
	if [ "$name" == "" ];
	then
		echo "$host not found";
		return
	fi

	$(expect -i -c '
		set host [lindex $argv 0]
		set pass [lindex $argv 1]
		spawn echo $host $pass
	' $user@$host $pass $cmd)
fi


if [ "$target" == "scp" ];
then
	return
fi


