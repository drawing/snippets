cs_selector = 0x10
ds_selector = 0x8

video = 0xb8000
# video_page = 0x4013000
video_page = 0x0000000

page_dir_base1	= 0x200000
page_tbl_base1	= 0x201000
page_dir_base2	= 0x400000
page_tbl_base2	= 0x401000
.code16
.text
	movw %cs, %ax
	movw %ax, %ds

	cli

	# enable A20
	inb $0x92, %al
	orb $0x02, %al
	outb %al, $0x92

	lgdt gdtr_data

	movl %cr0, %eax
	orb $0x01, %al
	movl %eax, %cr0

	ljmpl $cs_selector, $start_code32

gdtr_data:
	.short 0x18
	.int gdt_table

gdt_table:
	.quad 0x0
	.quad 0x00CF92000000FFFF
	.quad 0x00CF9A000000FFFF

.code32
.text
start_code32:
	movw $ds_selector, %ax
	movw %ax, %ds

	movl $page_dir_base1, %ebx
	movl $0x01, %eax
	#movl $0x201007, (%ebx, %eax, 4)
	movl $0x201007, (0x200004)
	movl $page_tbl_base1, %ebx
	movl $0x10, %eax
	movl $0x5011b, (%ebx, %eax, 4)
	movl $0x10, %eax
	#movl $0xb8007, (%ebx, %eax, 4)

	movl $0x201007, (0x200000)
	movl $0x201007, (0x200004)
	movl $0x201007, (0x200008)
	movl $0x201007, (0x20000C)
	movl $0x00007,  (0x201000)
	movl $0xb8007,  (0x201004)
	movl $0xb8007,  (0x201008)
	movl $0xb8007,  (0x20100C)

	movl $page_dir_base2, %ebx
	movl $0x2411b, 0x10(%ebx)
	movl $page_tbl_base2, %ebx
	movl $0x5111b, 0x12(%ebx)

	movl $'A', (0x50002);
	movl $'B', (0x51002);

	movl $page_dir_base1, %eax
	movl %eax, %cr3

	movl %cr0, %eax
	orl $0x80000000, %eax
	movl %eax, %cr0

	movl $0, %edi
	movb $0x02, %ah
	movb $'D', %al
	movw %ax, (0x001000)

	jmp .

.org 510, 0
	.short 0xAA55

