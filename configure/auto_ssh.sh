# usage auto_login.sh configure
# configure file:
#	servername	127.0.0.1	username	pass	command
#	server-8-1	localhost	cppbreak	1234	command
#
# add to system:
# ~/.bashrc
#	. auto_ssh.sh init configure
#
# usage:
# 	auto servername
#

# usage auto_rsync configure
# configure file:
# rsync_name	pass	/usr/local/	username@host:/usr/local/backup
# add to system:
# ~/.bashrc
# 	. auto_ssh.sh init_rsync configure
# usage:
# autorync rsync_name


script_name=auto_ssh.sh

[ $# -le 1 ] && { echo "Usage: $script_name command configure"; return;}

target=$1
config=$2

if [ "$target" == "init" ];
then
	alias auto="$script_name login $2";
	alias aucp="$script_name scp $2";

	name_list=`awk '{printf("%s ", $1)}' $config`;
	complete -W "$name_list" auto;
	return;
fi

if [ "$target" == "init_rsync" ];
then
	alias autorync="$script_name rsync $2";

	name_list=`awk '{printf("%s ", $1)}' $config`;
	complete -W "$name_list" autorync;
	return;
fi

if [ "$target" == "login" ];
then
	host_name=$3

	read name host user pass <<< $(awk -v key="$host_name" '{if ($1==key)print $1,$2,$3,$4};' $config);
	if [ "$name" == "" ];
	then
		echo "\"$host_name\" not found";
		exit 1;
	fi
	cmd=$(awk -v key="$host_name" '{if ($1==key){for (i=5;i<=NF;i++)printf($i" ")}};' $config);

	echo login: $host $user
	echo exec: $cmd

	tmux renamew $name
	expect_login.sh $host $user $pass "$cmd"
	echo exit: $name
	tmux renamew localhost
fi

if [ "$target" == "rsync" ];
then
	rsync_name=$3

	read name pass src dst <<< $(awk -v key="$rsync_name" '{if ($1==key)print $1,$2,$3,$4};' $config);
	if [ "$name" == "" ];
	then
		echo "\"$rsync_name\" not found";
		exit 1;
	fi

	echo rsync: "$src -> $dst"

	expect_rsync.sh $pass $src $dst
	echo finish: $name
fi

if [ "$target" == "scp" ];
then
	exit 0;
fi


