PROGRAM ?= app
PACKAGE_NAME ?= $(PROGRAM)
ROOT ?= ../..
ROOT_ABS := $(abspath $(ROOT))
STDLIB_DIR_ABS := $(ROOT_ABS)/stdlib

ABI_DIR = $(ROOT)/abi
MK_DIR = $(ROOT)/mk
BUILD_DIR = .build

GO = gccgo
GCC = gcc
ASM_COMPILER_FLAGS = -g -f elf32 -F dwarf
NASM = nasm $(ASM_COMPILER_FLAGS)
OBJCOPY = objcopy
SED = sed

ENTRYPOINT = go_0$(PACKAGE_NAME).Main
LDSCRIPT_TEMPLATE = $(MK_DIR)/static.lds.in
LDSCRIPT = $(BUILD_DIR)/$(PROGRAM).lds

GO_COMPILER_FLAGS = -m32 -c -nostdlib -nostdinc -fno-stack-protector -fno-split-stack -static -fno-leading-underscore -fno-common -fno-pie -g -ffunction-sections -fdata-sections -I. -I$(ROOT_ABS)
GCC_COMPILER_FLAGS = -m32 -c -ffunction-sections -fdata-sections -fno-pic -fno-pie -fno-stack-protector
LDFLAGS = -n -T $(LDSCRIPT) -m elf_i386 --no-ld-generated-unwind-info -z noexecstack -z relro -z now --gc-sections --entry=$(ENTRYPOINT)

APP_SOURCES = $(wildcard *.go)

PACKAGE_DIRS ?= kos ui
PACKAGE_OBJS =
PACKAGE_GOXS =
PREVIOUS_PACKAGE_GOXS =

define REGISTER_PACKAGE
PACKAGE_SOURCE_DIR_$(1) := $(if $(wildcard $(ROOT_ABS)/$(1)),$(ROOT_ABS)/$(1),$(if $(wildcard $(STDLIB_DIR_ABS)/$(1)),$(STDLIB_DIR_ABS)/$(1),$(error package source dir not found for $(1))))
PACKAGE_SOURCES_$(1) := $$(wildcard $$(PACKAGE_SOURCE_DIR_$(1))/*.go)
PACKAGE_SOURCE_FILES_$(1) := $$(notdir $$(PACKAGE_SOURCES_$(1)))

PACKAGE_OBJS += $(ROOT)/$(1).gccgo.o
PACKAGE_GOXS += $(ROOT)/$(1).gox

$(ROOT)/$(1).gccgo.o: $$(PACKAGE_SOURCES_$(1)) $(PREVIOUS_PACKAGE_GOXS)
	cd $$(PACKAGE_SOURCE_DIR_$(1)) && $(GO) $(GO_COMPILER_FLAGS) -o $(ROOT_ABS)/$(1).gccgo.o $$(PACKAGE_SOURCE_FILES_$(1))

$(ROOT)/$(1).gox: $(ROOT)/$(1).gccgo.o
	$(OBJCOPY) -j .go_export $$< $$@

PREVIOUS_PACKAGE_GOXS += $(ROOT)/$(1).gox
endef

$(foreach pkg,$(PACKAGE_DIRS),$(eval $(call REGISTER_PACKAGE,$(pkg))))

APP_OBJ = $(PROGRAM).gccgo.o

ABI_OBJS = $(ABI_DIR)/syscalls_i386.o $(ABI_DIR)/runtime_gccgo.o
OBJS = $(ABI_OBJS) $(PACKAGE_OBJS) $(APP_OBJ)
INTERMEDIATE_ARTIFACTS = $(ABI_OBJS) $(PACKAGE_OBJS) $(PACKAGE_GOXS) $(APP_OBJ) $(LDSCRIPT)

.PHONY: all clean link

all: $(PROGRAM).kex

clean:
	rm -rf $(BUILD_DIR)
	rm -f $(INTERMEDIATE_ARTIFACTS) $(PROGRAM).kex

link: $(PROGRAM).kex

$(BUILD_DIR):
	mkdir -p $@

$(LDSCRIPT): $(LDSCRIPT_TEMPLATE) | $(BUILD_DIR)
	$(SED) 's/@ENTRYPOINT@/$(ENTRYPOINT)/g' $< > $@

$(PROGRAM).kex: $(OBJS) $(PACKAGE_GOXS) $(LDSCRIPT)
	ld $(LDFLAGS) -o $(PROGRAM).kex $(OBJS)
	strip $(PROGRAM).kex
	$(OBJCOPY) $(PROGRAM).kex -O binary
	rm -f $(INTERMEDIATE_ARTIFACTS)
	rmdir $(BUILD_DIR) 2>/dev/null || true

$(APP_OBJ): $(APP_SOURCES) $(PACKAGE_GOXS)
	$(GO) $(GO_COMPILER_FLAGS) -o $@ $(APP_SOURCES)

$(ABI_DIR)/runtime_gccgo.o: $(ABI_DIR)/runtime_gccgo.c
	$(GCC) $(GCC_COMPILER_FLAGS) $< -o $@

$(ABI_DIR)/syscalls_i386.o: $(ABI_DIR)/syscalls_i386.asm
	$(NASM) $< -o $@
