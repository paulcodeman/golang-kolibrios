.PHONY: all example hello strings sysinfo message clean clean-example clean-hello clean-strings clean-sysinfo clean-message rebuild-example rebuild-hello rebuild-strings rebuild-sysinfo rebuild-message rebuild-all

all: example hello strings sysinfo message

example:
	$(MAKE) -C cmd/example all

hello:
	$(MAKE) -C cmd/hello all

strings:
	$(MAKE) -C cmd/strings all

sysinfo:
	$(MAKE) -C cmd/sysinfo all

message:
	$(MAKE) -C cmd/message all

clean: clean-example clean-hello clean-strings clean-sysinfo clean-message

clean-example:
	$(MAKE) -C cmd/example clean

clean-hello:
	$(MAKE) -C cmd/hello clean

clean-strings:
	$(MAKE) -C cmd/strings clean

clean-sysinfo:
	$(MAKE) -C cmd/sysinfo clean

clean-message:
	$(MAKE) -C cmd/message clean

rebuild-example:
	$(MAKE) -C cmd/example clean all

rebuild-hello:
	$(MAKE) -C cmd/hello clean all

rebuild-strings:
	$(MAKE) -C cmd/strings clean all

rebuild-sysinfo:
	$(MAKE) -C cmd/sysinfo clean all

rebuild-message:
	$(MAKE) -C cmd/message clean all

rebuild-all:
	$(MAKE) clean all
