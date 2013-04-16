# Copyright 2013 Prometheus Team
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.SUFFIXES:

include ../Makefile.INCLUDE

all: dependencies-stamp

bison-stamp: bison-implementation-$(UNAME)-stamp
	[ -x "$$(which bison)" ] || { echo "bison not found." ; false ; }
	touch $@

bison-implementation-Darwin-stamp:
	[ -x "$$(which bison)" ] || $(BREW_INSTALL) bison
	touch $@

bison-implementation-Linux-stamp:
	[ -x "$$(which bison)" ] || $(APT_GET_INSTALL) bison
	touch $@

cache-stamp:
	$(MAKE) -C cache
	touch $@

cc-stamp: cc-implementation-$(UNAME)-stamp
	[ -x "$$(which cc)" ] || { echo "cc not found." ; false ; }
	touch $@

cc-implementation-Darwin-stamp:
	[ -x "$$(which cc)" ] || { echo "Install XCode?" ; false ; }
	touch $@

cc-implementation-Linux-stamp:
	[ -x "$$(which cc)" ] || $(APT_GET_INSTALL) build-essential
	touch $@

dependencies-stamp: cache-stamp go-stamp goprotobuf-stamp gorest-stamp goskiplist-stamp instrumentation-stamp levigo-stamp
	touch $@

go-stamp: gvm-stamp
	[ -x "$$(which go)" ] || { echo "go not found." ; false ; }
	touch $@

goprotobuf-stamp: go-stamp protoc-stamp
	$(GO_GET) code.google.com/p/goprotobuf/proto
	$(GO_GET) code.google.com/p/goprotobuf/protoc-gen-go
	touch $@

gorest-stamp: go-stamp
	$(GO_GET) code.google.com/p/gorest
	touch $@

goskiplist-stamp: go-stamp
	$(GO_GET) github.com/ryszard/goskiplist/skiplist
	touch $@

gperftools-stamp: cache-stamp cache/gperftools-$(GPERFTOOLS_VERSION).tar.gz cc-stamp stack-unwind-support-stamp
	tar xzvf cache/gperftools-$(GPERFTOOLS_VERSION).tar.gz -C dirty
	cd dirty/gperftools-$(GPERFTOOLS_VERSION) && ./configure --prefix=$(PREFIX)
	make -C dirty/gperftools-$(GPERFTOOLS_VERSION)
	make -C dirty/gperftools-$(GPERFTOOLS_VERSION) install
	touch $@

gvm-stamp: cc-stamp
	[ -x "$$(which gvm)" ] || { echo "gvm not found; install manually from https://github.com/moovweb/gvm." ; false ; }
	-gvm install go$(GO_VERSION)
	#
	# ********** **********
	# gvm installed go$(GO_VERSION) if it did not exist already.
  #
  # Prometheus expects go$(GO_VERSION) to be used for building.
	# This can be achieved by running ``gvm use go$(GO_VERSION)`` before any
	# other Prometheus make tools.
	# ********** **********
	#
	[ -z "${CI}" ] && read -p "Press [ENTER] upon reading." _
	touch $@

instrumentation-stamp: go-stamp
	$(GO_GET) github.com/prometheus/client_golang
	touch $@

leveldb-stamp: cache-stamp cache/leveldb-$(LEVELDB_VERSION).tar.gz cc-stamp gperftools-stamp rsync-stamp snappy-stamp
	tar xzvf cache/leveldb-$(LEVELDB_VERSION).tar.gz -C dirty
	# The test that LevelDB uses to test for TCMalloc and Snappy is naive and
	# does not respect LDFLAGS.  :-(
	CXXFLAGS="-lsnappy -ltcmalloc $(LDFLAGS)" LDFLAGS="-lsnappy -ltcmalloc $(LDFLAGS)" $(MAKE) -C dirty/leveldb-$(LEVELDB_VERSION)
	rsync -av "dirty/leveldb-$(LEVELDB_VERSION)/include/" "$(PREFIX)/include/"
	-[ "$(UNAME)" = "Linux" ] && { rsync -av "dirty/leveldb-$(LEVELDB_VERSION)/"*.*so* "$(PREFIX)/lib/" ; }
	-[ "$(UNAME)" = "Darwin" ] && { rsync -av "dirty/leveldb-$(LEVELDB_VERSION)/"*.*dylib* "$(PREFIX)/lib/" ; }
	rsync -av "dirty/leveldb-$(LEVELDB_VERSION)/"*.a "$(PREFIX)/lib/"
	touch $@

levigo-stamp: cc-stamp gperftools-stamp go-stamp leveldb-stamp snappy-stamp
	$(GO_GET) github.com/jmhodges/levigo
	touch $@

libunwind-stamp:
	$(APT_GET_INSTALL) libunwind7
	$(APT_GET_INSTALL) libunwind7-dev
	touch $@

mercurial-implementation-Darwin-stamp:
	[ -x "$$(which hg)" ] || $(BREW_INSTALL) mercurial
	touch $@

mercurial-implementation-Linux-stamp:
	[ -x "$$(which hg)" ] || $(APT_GET_INSTALL) mercurial
	touch $@

mercurial-stamp: mercurial-implementation-$(UNAME)-stamp
	[ -x "$$(which hg)" ] || { echo "hg not found." ; false ; }
	touch $@

noop-target-stamp:
	echo "Not doing anything."
	touch $@

protoc-stamp: cache-stamp cache/protobuf-$(PROTOCOL_BUFFERS_VERSION).tar.bz2 cc-stamp
	tar xjvf cache/protobuf-$(PROTOCOL_BUFFERS_VERSION).tar.bz2 -C dirty
	cd dirty/protobuf-$(PROTOCOL_BUFFERS_VERSION) && ./configure --prefix="$(PREFIX)"
	$(MAKE) -C dirty/protobuf-$(PROTOCOL_BUFFERS_VERSION)
	$(MAKE) -C dirty/protobuf-$(PROTOCOL_BUFFERS_VERSION) install
	[ -x "$$(which protoc)" ] || { echo "protoc not found." ; false ; }
	touch $@

rsync-implementation-Darwin-stamp:
	[ -x "$$(which rsync)" ] || $(BREW_INSTALL) rsync
	touch $@

rsync-implementation-Linux-stamp:
	[ -x "$$(which rsync)" ] || $(APT_GET_INSTALL) rsync

rsync-stamp: rsync-implementation-$(UNAME)-stamp
	[ -x "$$(which rsync)" ] || { echo "rsync not found." ; false ; }
	touch $@

snappy-stamp: cache-stamp cache/snappy-$(SNAPPY_VERSION).tar.gz cc-stamp gperftools-stamp
	tar xzvf cache/snappy-$(SNAPPY_VERSION).tar.gz -C dirty
	cd dirty/snappy-$(SNAPPY_VERSION) && ./configure --prefix="$(PREFIX)"
	$(MAKE) -C dirty/snappy-$(SNAPPY_VERSION)
	$(MAKE) -C dirty/snappy-$(SNAPPY_VERSION) install
	touch $@

ifeq ($(UNAME), Linux)
stack-unwind-support-stamp: libunwind-stamp
	touch $@
else
stack-unwind-support-stamp: noop-target-stamp
	touch $@
endif

vim-implementation-Darwin-stamp:
	[ -x "$$(which vim)" ] || $(BREW_INSTALL) vim
	touch $@

vim-implementation-Linux-stamp:
	[ -x "$$(which vim)" ] || $(APT_GET_INSTALL) vim
	touch $@

vim-stamp: vim-implementation-$(UNAME)-stamp
	touch $@

# go-dependencies-stamp: go-stamp gorest goskiplist instrumentation levigo mercurial protoc
# 	touch $@

clean:
	$(MAKE) -C cache clean
	$(MAKE) -C dirty clean
	$(MAKE) -C root clean
	rm -rf *-stamp


.PHONY: clean