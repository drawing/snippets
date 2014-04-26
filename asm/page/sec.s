cs_selector = 0x10
ds_selector = 0x8

video = 0xb8000

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
	movw %ax, %gs

	movb $0x02, %ah
	movb $'D', %al
	movw %ax, %gs:(video)
	jmp .

.org 510, 0
	.short 0xAA55

