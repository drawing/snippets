set -x
program=$1
as $program -o a.o
ld a.o -Tbios.lds
objcopy -R .pdr -R .comment -R .note -S -O binary a.out disk.img
