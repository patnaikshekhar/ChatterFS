package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hanwen/go-fuse/fuse"

	"github.com/hanwen/go-fuse/fuse/nodefs"

	"github.com/hanwen/go-fuse/fuse/pathfs"
)

// ChatterFs is the root filesystem
type ChatterFs struct {
	pathfs.FileSystem
}

// GetAttr gets the attributes for files
func (cfs *ChatterFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	switch name {
	case "file.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644,
			Size: uint64(len(name)),
		}, fuse.OK
	case "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	}

	return nil, fuse.ENOENT
}

// OpenDir opens a directory
func (cfs *ChatterFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	if name == "" {
		c = []fuse.DirEntry{{Name: "file.txt", Mode: fuse.S_IFREG}}
		return c, fuse.OK
	}

	return nil, fuse.ENOENT
}

// Open opens a file
func (cfs *ChatterFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	if name != "file.txt" {
		return nil, fuse.ENOENT
	} else {
		if flags&fuse.O_ANYWRITE != 0 {
			return nil, fuse.EPERM
		}

		return nodefs.NewDataFile([]byte(name)), fuse.OK
	}
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Fatal("Usage:\n chatterFS MOUNTPOINT")
	}

	fmt.Println("file", flag.Arg(0))

	nfs := pathfs.NewPathNodeFs(&ChatterFs{FileSystem: pathfs.NewDefaultFileSystem()}, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail %v\n", err)
	}

	server.Serve()
}
