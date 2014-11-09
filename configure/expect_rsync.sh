#!/usr/bin/expect --

if { [llength $argv] < 3} {
	puts "usage: $argv0 password src dst"
	exit 1
}

set timeout 45
set password [lindex $argv 0]
set src [lindex $argv 1]
set dst [lindex $argv 2]

trap {
	set rows [stty rows]
	set cols [stty columns]
	stty rows $rows columns $cols < $spawn_out(slave,name)
} WINCH

spawn rsync -avH ssh $src $dst

expect {
	"password:" {
		send "$password\r"
	}

	"(yes/no)?" {
		send "yes\r"
		expect "password:" {
			send "$password\r"
		}
	}
	timeout {
		exit 1
	}
}

interact

