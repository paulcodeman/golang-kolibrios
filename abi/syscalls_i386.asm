SECTION .text

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
    mov eax, 5
    mov ebx, [ebp+8]
    int 0x40
    mov esp, ebp
    pop ebp
    ret

go_0kos.Event:
    mov eax, 10
    int 0x40
    ret

go_0kos.GetButtonID:
    mov eax, 17
    int 0x40
    test al, al
    jnz .no_button
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
    mov eax, 12
    mov ebx, [ebp+8]
    int 0x40
    mov esp, ebp
    pop ebp
    ret

go_0kos.Window:
    push ebp
    mov ebp, esp
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
    mov edi, [ebp+24]
    xor eax, eax
    int 0x40
    mov esp, ebp
    pop ebp
    ret

go_0kos.WriteText:
    push ebp
    mov ebp, esp
    mov eax, 4
    mov ebx, [ebp+8]
    shl ebx, 16
    mov bx, [ebp+12]
    mov ecx, [ebp+16]
    mov edx, [ebp+20]
    mov esi, [ebp+24]
    int 0x40
    mov esp, ebp
    pop ebp
    ret

go_0kos.DrawLine:
    push ebp
    mov ebp, esp
    mov ebx, [ebp+8]
    shl ebx, 16
    mov bx, [ebp+16]
    mov ecx, [ebp+12]
    shl ecx, 16
    mov cx, [ebp+20]
    mov edx, [ebp+24]
    mov eax, 38
    int 0x40
    mov esp, ebp
    pop ebp
    ret

go_0kos.DrawBar:
    push ebp
    mov ebp, esp
    mov eax, 13
    mov ebx, [ebp+8]
    shl ebx, 16
    mov bx, [ebp+16]
    mov ecx, [ebp+12]
    shl ecx, 16
    mov cx, [ebp+20]
    mov edx, [ebp+24]
    int 0x40
    mov esp, ebp
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
    mov edx, [esp+4]
    mov eax, 63
    mov ebx, 1
.next_char:
    mov cl, [edx]
    test cl, cl
    jz .done
    int 0x40
    inc edx
    jmp .next_char
.done:
    ret

go_0kos.CreateButton:
    push ebp
    mov ebp, esp
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
    mov esp, ebp
    pop ebp
    ret

global malloc
global free
global realloc

malloc:
    mov eax, 68
    mov ebx, 12
    mov ecx, [esp+4]
    int 0x40
    ret

free:
    mov eax, 68
    mov ebx, 13
    mov ecx, [esp+4]
    int 0x40
    ret

realloc:
    mov eax, 68
    mov ebx, 13
    mov edx, [esp+4]
    mov ecx, [esp+8]
    int 0x40
    ret

SECTION .data

__hexdigits:
    db '0123456789ABCDEF'
