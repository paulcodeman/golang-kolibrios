SECTIONS
{
    . = 0x00000;

    .text :
    {
        LONG(0x554E454D);
        LONG(0x31305445);
        LONG(1);
        LONG(go_0example.Main);
        LONG(__end);
        LONG(0x10000);
        LONG(0x10000);
        LONG(0);
        LONG(0);

        *(.text)
    }

    .rodata :
    {
        *(.rodata)
        *(.rodata.*)
    }

    .data :
    {
        *(.data)
    }

    .bss :
    {
        *(COMMON)
        *(.bss)
    }

    .eh_frame :
    {
        *(.eh_frame)
    }

    .group :
    {
        *(.group)
    }

    __end = .;

    .note.GNU-stack : { }
}
