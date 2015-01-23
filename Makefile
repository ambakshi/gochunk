.PHONY: all sanity

UNAME:=$(shell uname -s)

ifeq ($(UNAME),Darwin)
MD5=md5
else
MD5=md5sum -b
endif

all: gochunk

gochunk: main.go chunk.go
	go build


sanity: gochunk
	./$^ c $^ | tee $^.manifest
	@for chunk in `cat $^.manifest | awk '{print $$1}'`; do \
        fname="`echo $$chunk | cut -c1-2`/`echo $$chunk | cut -c3-4`/`echo $$chunk | cut -c5-6`/$$chunk"; \
        gunzip -c < chunks/$$fname ; done > $^.chk
	$(MD5) $^ $^.chk
	rm -f $^.chk

clean:
	rm -f gochunk
	rm -rf ./chunks/
