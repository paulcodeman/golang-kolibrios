SECTION .text

extern runtime_prepare_window_title

global go_0kos.Sleep
global go_0kos.Event
global go_0kos.GetButtonID
global go_0kos.CreateButton
global go_0kos.Exit
global go_0kos.Redraw
global go_0kos.Window
global go_0kos.WriteText
global go_0kos.GetTime
global go_0kos.DrawLine
global go_0kos.DrawBar
global go_0kos.DebugOutHex
global go_0kos.DebugOutChar
global go_0kos.DebugOutStr

go_0kos.Sleep:
    push ebp
    mov ebp, esp
    push ebx
    mov eax, 5
    mov ebx, [ebp+8]
    int 0x40
    pop ebx
    pop ebp
    ret

go_0kos.Event:
    mov eax, 10
    int 0x40
    ret

go_0kos.GetButtonID:
    mov eax, 17
    int 0x40
    cmp eax, 1
    je .no_button
    shr eax, 8
    ret
.no_button:
    xor eax, eax
    dec eax
    ret

go_0kos.Exit:
    mov eax, -1
    int 0x40
    ret

go_0kos.Redraw:
    push ebp
    mov ebp, esp
    push ebx
    mov eax, 12
    mov ebx, [ebp+8]
    int 0x40
    pop ebx
    pop ebp
    ret

go_0kos.Window:
    push ebp
    mov ebp, esp
    push ebx
    push esi
    push edi
    push dword [ebp+28]
    push dword [ebp+24]
    call runtime_prepare_window_title
    add esp, 8
    mov edi, eax
    mov ebx, [ebp+8]
    shl ebx, 16
    or ebx, [ebp+16]
    mov ecx, [ebp+12]
    shl ecx, 16
    or ecx, [ebp+20]
    mov edx, 0x14
    shl edx, 24
    or edx, 0xFFFFFF
    mov esi, 0x808899FF
    xor eax, eax
    int 0x40
    pop edi
    pop esi
    pop ebx
    pop ebp
    ret

go_0kos.WriteText:
    push ebp
    mov ebp, esp
    push ebx
    push esi
    mov eax, 4
    mov ebx, [ebp+8]
    shl ebx, 16
    mov bx, [ebp+12]
    mov ecx, [ebp+16]
    and ecx, 0x00FFFFFF
    or ecx, 0x30000000
    mov edx, [ebp+20]
    mov esi, [ebp+24]
    int 0x40
    pop esi
    pop ebx
    pop ebp
    ret

go_0kos.DrawLine:
    push ebp
    mov ebp, esp
    push ebx
    mov ebx, [ebp+8]
    shl ebx, 16
    mov bx, [ebp+16]
    mov ecx, [ebp+12]
    shl ecx, 16
    mov cx, [ebp+20]
    mov edx, [ebp+24]
    mov eax, 38
    int 0x40
    pop ebx
    pop ebp
    ret

go_0kos.DrawBar:
    push ebp
    mov ebp, esp
    push ebx
    mov eax, 13
    mov ebx, [ebp+8]
    shl ebx, 16
    mov bx, [ebp+16]
    mov ecx, [ebp+12]
    shl ecx, 16
    mov cx, [ebp+20]
    mov edx, [ebp+24]
    int 0x40
    pop ebx
    pop ebp
    ret

go_0kos.GetTime:
    mov eax, 3
    int 0x40
    ret

go_0kos.DebugOutHex:
    mov eax, [esp+4]
    mov edx, 8
.next_hex_digit:
    rol eax, 4
    movzx ecx, al
    and cl, 0x0F
    mov cl, [__hexdigits + ecx]
    pushad
    mov eax, 63
    mov ebx, 1
    int 0x40
    popad
    dec edx
    jnz .next_hex_digit
    ret

go_0kos.DebugOutChar:
    mov al, [esp+4]
    pushf
    pushad
    mov cl, al
    mov eax, 63
    mov ebx, 1
    int 0x40
    popad
    popf
    ret

go_0kos.DebugOutStr:
    push ebx
    push esi
    mov edx, [esp+12]
    mov esi, [esp+16]
    mov eax, 63
    mov ebx, 1
.next_char:
    test esi, esi
    jz .done
    mov cl, [edx]
    int 0x40
    inc edx
    dec esi
    jmp .next_char
.done:
    pop esi
    pop ebx
    ret

go_0kos.CreateButton:
    push ebp
    mov ebp, esp
    push ebx
    push esi
    mov eax, 8
    mov ebx, [ebp+8]
    shl ebx, 16
    mov bx, [ebp+16]
    mov ecx, [ebp+12]
    shl ecx, 16
    mov cx, [ebp+20]
    mov edx, [ebp+24]
    mov esi, [ebp+28]
    int 0x40
    pop esi
    pop ebx
    pop ebp
    ret

global malloc
global free
global realloc

malloc:
    push ebx
    call __ensure_heap_initialized
    test eax, eax
    jz .malloc_failed
    mov eax, 68
    mov ebx, 12
    mov ecx, [esp+8]
    int 0x40
    pop ebx
    ret
.malloc_failed:
    xor eax, eax
    pop ebx
    ret

free:
    push ebx
    call __ensure_heap_initialized
    test eax, eax
    jz .free_failed
    mov eax, 68
    mov ebx, 13
    mov ecx, [esp+8]
    int 0x40
    pop ebx
    ret
.free_failed:
    xor eax, eax
    pop ebx
    ret

realloc:
    push ebx
    call __ensure_heap_initialized
    test eax, eax
    jz .realloc_failed
    mov eax, 68
    mov ebx, 20
    mov edx, [esp+8]
    mov ecx, [esp+12]
    int 0x40
    pop ebx
    ret
.realloc_failed:
    xor eax, eax
    pop ebx
    ret

__ensure_heap_initialized:
    cmp dword [__heap_initialized], 0
    jne .ready
    mov eax, 68
    mov ebx, 11
    int 0x40
    test eax, eax
    jz .failed
    mov dword [__heap_initialized], 1
.ready:
    mov eax, 1
    ret
.failed:
    xor eax, eax
    ret

SECTION .data

__hexdigits:
    db '0123456789ABCDEF'

SECTION .bss

__heap_initialized:
    resd 1
