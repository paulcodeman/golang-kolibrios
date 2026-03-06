.PHONY: all check-runtime example hello strings slices interfaces emptyiface sysinfo message ipc clean clean-example clean-hello clean-strings clean-slices clean-interfaces clean-emptyiface clean-sysinfo clean-message clean-ipc rebuild-example rebuild-hello rebuild-strings rebuild-slices rebuild-interfaces rebuild-emptyiface rebuild-sysinfo rebuild-message rebuild-ipc rebuild-all

all: example hello strings slices interfaces emptyiface sysinfo message ipc

check-runtime:
	bash ./scripts/check-runtime-probes.sh

example:
	$(MAKE) -C cmd/example all

hello:
	$(MAKE) -C cmd/hello all

strings:
	$(MAKE) -C cmd/strings all

slices:
	$(MAKE) -C cmd/slices all

interfaces:
	$(MAKE) -C cmd/interfaces all

emptyiface:
	$(MAKE) -C cmd/emptyiface all

sysinfo:
	$(MAKE) -C cmd/sysinfo all

message:
	$(MAKE) -C cmd/message all

ipc:
	$(MAKE) -C cmd/ipc all

clean: clean-example clean-hello clean-strings clean-slices clean-interfaces clean-emptyiface clean-sysinfo clean-message clean-ipc

clean-example:
	$(MAKE) -C cmd/example clean

clean-hello:
	$(MAKE) -C cmd/hello clean

clean-strings:
	$(MAKE) -C cmd/strings clean

clean-slices:
	$(MAKE) -C cmd/slices clean

clean-interfaces:
	$(MAKE) -C cmd/interfaces clean

clean-emptyiface:
	$(MAKE) -C cmd/emptyiface clean

clean-sysinfo:
	$(MAKE) -C cmd/sysinfo clean

clean-message:
	$(MAKE) -C cmd/message clean

clean-ipc:
	$(MAKE) -C cmd/ipc clean

rebuild-example:
	$(MAKE) -C cmd/example clean all

rebuild-hello:
	$(MAKE) -C cmd/hello clean all

rebuild-strings:
	$(MAKE) -C cmd/strings clean all

rebuild-slices:
	$(MAKE) -C cmd/slices clean all

rebuild-interfaces:
	$(MAKE) -C cmd/interfaces clean all

rebuild-emptyiface:
	$(MAKE) -C cmd/emptyiface clean all

rebuild-sysinfo:
	$(MAKE) -C cmd/sysinfo clean all

rebuild-message:
	$(MAKE) -C cmd/message clean all

rebuild-ipc:
	$(MAKE) -C cmd/ipc clean all

rebuild-all:
	$(MAKE) clean all
