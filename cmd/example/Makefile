PROGRAM = example
ROOT = ../..
ABI_DIR = $(ROOT)/abi
KOS_DIR = $(ROOT)/kos
UI_DIR = $(ROOT)/ui

GO = gccgo
GCC = gcc
ASM_COMPILER_FLAGS = -g -f elf32 -F dwarf
NASM = nasm $(ASM_COMPILER_FLAGS)
OBJCOPY = objcopy

GO_COMPILER_FLAGS = -m32 -c -nostdlib -nostdinc -fno-stack-protector -fno-split-stack -static -fno-leading-underscore -fno-common -fno-pie -g -ffunction-sections -fdata-sections -I.
GCC_COMPILER_FLAGS = -m32 -c -ffunction-sections -fdata-sections
LDFLAGS = -n -T static.lds -m elf_i386 --no-ld-generated-unwind-info -z noexecstack -z relro -z now --gc-sections --entry=go_0example.Main

KOS_SOURCES = $(wildcard $(KOS_DIR)/*.go)
UI_SOURCES = $(wildcard $(UI_DIR)/*.go)
APP_SOURCES = $(wildcard *.go)

KOS_OBJ = $(ROOT)/kos.gccgo.o
KOS_GOX = $(ROOT)/kos.gox
UI_OBJ = $(ROOT)/ui.gccgo.o
UI_GOX = $(ROOT)/ui.gox
APP_OBJ = $(PROGRAM).gccgo.o

ABI_OBJS = $(ABI_DIR)/syscalls_i386.o $(ABI_DIR)/runtime_gccgo.o
OBJS = $(ABI_OBJS) $(KOS_OBJ) $(UI_OBJ) $(APP_OBJ)

.PHONY: all clean link

all: $(OBJS) $(KOS_GOX) $(UI_GOX) link

clean:
	rm -f $(ABI_DIR)/*.o $(ROOT)/*.gccgo.o $(ROOT)/*.gox $(APP_OBJ) $(PROGRAM).kex

link:
	ld $(LDFLAGS) -o $(PROGRAM).kex $(OBJS)
	strip $(PROGRAM).kex
	$(OBJCOPY) $(PROGRAM).kex -O binary

$(KOS_OBJ): $(KOS_SOURCES)
	$(GO) $(GO_COMPILER_FLAGS) -o $@ $(KOS_SOURCES)

$(KOS_GOX): $(KOS_OBJ)
	$(OBJCOPY) -j .go_export $< $@

$(UI_OBJ): $(UI_SOURCES) $(KOS_GOX)
	$(GO) $(GO_COMPILER_FLAGS) -o $@ $(UI_SOURCES)

$(UI_GOX): $(UI_OBJ)
	$(OBJCOPY) -j .go_export $< $@

$(APP_OBJ): $(APP_SOURCES) $(KOS_GOX) $(UI_GOX)
	$(GO) $(GO_COMPILER_FLAGS) -o $@ $(APP_SOURCES)

$(ABI_DIR)/runtime_gccgo.o: $(ABI_DIR)/runtime_gccgo.c
	$(GCC) $(GCC_COMPILER_FLAGS) $< -o $@

$(ABI_DIR)/syscalls_i386.o: $(ABI_DIR)/syscalls_i386.asm
	$(NASM) $< -o $@
