; vim: set ft=nasm :

section .bss
align 4

; Reserve 16K for our stack. Stacks should be aligned to 16 byte boundaries.
stack_bottom:
	resb 16384 	; 16 KiB
stack_top:

; According to the "ELF handling for TLS" document section 4.3.2
; (https://www.akkadia.org/drepper/tls.pdf) for the GNU variant of the IA-32 ABI, 
; gs:0x00 contains a pointer to the TCB. Variables in the TLS are stored 
; before the TCB and are accessed using negative offsets from the TCB address.
g0_ptr:	        resd 1 
tcb_ptr:        resd 1 

section .text
bits 32
align 4

MULTIBOOT_MAGIC equ 0x36d76289

G_STACK_LO equ 0x0
G_STACK_HI equ 0x4
G_STACKGUARD0 equ 0x8

err_unsupported_bootloader db '[rt0] kernel not loaded by multiboot-compliant bootloader', 0

;------------------------------------------------------------------------------
; Kernel arch-specific entry point
;
; The boot loader will jump to this symbol after setting up the CPU according
; to the multiboot standard. At this point:
; - A20 is enabled
; - The CPU is using 32-bit protected mode
; - Interrupts are disabled
; - Paging is disabled
; - EAX contains the magic value ‘0x36d76289’; the presence of this value indicates
;   to the operating system that it was loaded by a Multiboot-compliant boot loader
; - EBX contains the 32-bit physical address of the Multiboot information structure
;------------------------------------------------------------------------------
global _rt0_entry
_rt0_entry:
	cmp eax, MULTIBOOT_MAGIC
	jne unsupported_bootloader

	; Initalize our stack by pointing ESP to the BSS-allocated stack. In x86,
	; stack grows downwards so we need to point ESP to stack_top
	mov esp, stack_top

	; Enable SSE/AVX
	call _rt0_enable_sse

 	; Load initial GDT
 	call _rt0_load_gdt

	; init g0 so we can invoke Go functions. For now we use hardcoded offsets 
	; that correspond to the g struct definition in src/runtime/runtime2.go
	extern runtime.g0
	mov dword [runtime.g0 + G_STACK_LO], stack_bottom
	mov dword [runtime.g0 + G_STACK_HI], stack_top
	mov dword [runtime.g0 + G_STACKGUARD0], stack_bottom
	mov dword [g0_ptr], runtime.g0

	; jump into the go code
	extern main.main
	call main.main

	; Main should never return; halt the CPU
halt:
	cli
	hlt

unsupported_bootloader:
	mov edi, err_unsupported_bootloader
	call write_string
	jmp halt
.end:

;------------------------------------------------------------------------------
; Write the NULL-terminated string contained in edi to the screen using white
; text on red background.  Assumes that text-mode is enabled and that its
; physical address is 0xb8000.
;------------------------------------------------------------------------------
write_string:
	push eax
	push ebx

	mov ebx,0xb8000
	mov ah, 0x4F
next_char:
	mov al, byte[edi]
	test al, al
	jz done

	mov word [ebx], ax
	add ebx, 2
	inc edi
	jmp next_char

done:
	pop ebx
	pop eax
	ret

;------------------------------------------------------------------------------
; Load GDT and flush CPU caches
;------------------------------------------------------------------------------

_rt0_load_gdt:
	push eax
	push ebx

	; Store the address to the TCB in tcb_ptr
	; and set up gs base address to it
	mov eax, tcb_ptr
	mov [tcb_ptr], eax
	mov ebx, gdt0_gs_seg
	mov [ebx+2], al
	mov [ebx+3], ah
	shr eax, 16
	mov [ebx+4], al

	lgdt [gdt0_desc]

	; GDT has been loaded but the CPU still has the previous GDT data in cache.
	; We need to manually update the descriptors and use a JMP command to set
	; the CS segment descriptor
	jmp CS_SEG:update_descriptors
update_descriptors:
	mov ax, DS_SEG
	mov ds, ax
	mov es, ax
	mov fs, ax
	mov ss, ax
	mov ax, GS_SEG
	mov gs, ax

	pop ebx
	pop eax
	ret

;------------------------------------------------------------------------------
; GDT definition
;------------------------------------------------------------------------------
%include "gdt.inc"

align 2
gdt0:

gdt0_nil_seg: GDT_ENTRY_32 0x00, 0x0, 0x0, 0x0				        ; nil descriptor (not used by CPU but required by some emulators)
gdt0_cs_seg:  GDT_ENTRY_32 0x00, 0xFFFFF, SEG_EXEC | SEG_R, SEG_GRAN_4K_PAGE    ; code descriptor
gdt0_ds_seg:  GDT_ENTRY_32 0x00, 0xFFFFF, SEG_NOEXEC | SEG_W, SEG_GRAN_4K_PAGE  ; data descriptor
gdt0_gs_seg:  GDT_ENTRY_32 0x00, 0xFFFFF, SEG_NOEXEC | SEG_W, SEG_GRAN_BYTE        ; TLS descriptor (required in order to use go segmented stacks)

gdt0_desc:
	dw gdt0_desc - gdt0 - 1  ; gdt size should be 1 byte less than actual length
	dd gdt0

NULL_SEG equ gdt0_nil_seg - gdt0
CS_SEG   equ gdt0_cs_seg - gdt0
DS_SEG   equ gdt0_ds_seg - gdt0
GS_SEG   equ gdt0_gs_seg - gdt0

;------------------------------------------------------------------------------
; Enable SSE support. Code taken from:
; http://wiki.osdev.org/SSE#Checking_for_SSE
;------------------------------------------------------------------------------
_rt0_enable_sse:
	push eax

	; check for SSE
	mov eax, 0x1
	cpuid
	test edx, 1<<25
	jz .no_sse

	; enable SSE
	mov eax, cr0
	and ax, 0xFFFB      ; clear coprocessor emulation CR0.EM
	or ax, 0x2          ; set coprocessor monitoring  CR0.MP
	mov cr0, eax
	mov eax, cr4
	or ax, 3 << 9       ; set CR4.OSFXSR and CR4.OSXMMEXCPT at the same time
	mov cr4, eax

	pop eax
	ret
.no_sse:
	cli
	hlt
