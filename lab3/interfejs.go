package main

import (
	"errors"
	"fmt"
	"time"
)

// Iterfejs definiujący obiekt w systemie plików
type FileSystemItem interface {
	Name() string
	Path() string
	Size() int64
	CreatedAt() time.Time
	ModifiedAt() time.Time
}

// Interfejs definiujący obiekty które mogą być odczttywane
type Readable interface {
	Read(p []byte) (n int, err error)
}

// Interfejs definiujący obiekty w których można dokonywać zapisu
type Writable interface {
	Write(p []byte) (n int, err error)
}

// Katalog definiuje pliki i podkatalogi
type Directory interface {
	FileSystemItem
	AddItem(item FileSystemItem) error
	RemoveItem(name string) error
	Items() []FileSystemItem
}

// Przykładowe komunikaty błędów, które można użyć
var (
	ErrItemExists       = errors.New("item already exists")
	ErrItemNotFound     = errors.New("item not found")
	ErrNotImplemented   = errors.New("operation not implemented")
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotDirectory     = errors.New("not a directory")
	ErrIsDirectory      = errors.New("is a directory")
)

type Plik struct {
	name       string
	path       string
	size       int64
	createdAt  time.Time
	modifiedAt time.Time
	data       []byte
}

func (f *Plik) Name() string          { return f.name }
func (f *Plik) Path() string          { return f.path }
func (f *Plik) Size() int64           { return int64(len(f.data)) }
func (f *Plik) CreatedAt() time.Time  { return f.createdAt }
func (f *Plik) ModifiedAt() time.Time { return f.modifiedAt }
func (f *Plik) Read(p []byte) (int, error) {
	copy(p, f.data)
	return len(f.data), nil
}
func (f *Plik) Write(p []byte) (int, error) {
	f.data = append(f.data, p...)
	f.modifiedAt = time.Now()
	return len(p), nil
}

type Katalog struct {
	name       string
	path       string
	size       int64
	createdAt  time.Time
	modifiedAt time.Time
	items      map[string]FileSystemItem
}

func (d *Katalog) Name() string          { return d.name }
func (d *Katalog) Path() string          { return d.path }
func (d *Katalog) Size() int64           { return d.size }
func (d *Katalog) CreatedAt() time.Time  { return d.createdAt }
func (d *Katalog) ModifiedAt() time.Time { return d.modifiedAt }
func (d *Katalog) AddItem(item FileSystemItem) error {
	if _, exists := d.items[item.Name()]; exists {
		return ErrItemExists
	}
	d.items[item.Name()] = item
	d.modifiedAt = time.Now()
	return nil
}
func (d *Katalog) RemoveItem(name string) error {
	if _, exists := d.items[name]; !exists {
		return ErrItemNotFound
	}
	delete(d.items, name)
	d.modifiedAt = time.Now()
	return nil
}
func (d *Katalog) Items() []FileSystemItem {
	var result []FileSystemItem
	for _, item := range d.items {
		result = append(result, item)
	}
	return result
}

type SymLink struct {
	name       string
	path       string
	size       int64
	createdAt  time.Time
	modifiedAt time.Time
	target     FileSystemItem
}

func (s *SymLink) Name() string          { return s.name }
func (s *SymLink) Path() string          { return s.path }
func (s *SymLink) Size() int64           { return s.size }
func (s *SymLink) CreatedAt() time.Time  { return s.createdAt }
func (s *SymLink) ModifiedAt() time.Time { return s.modifiedAt }

type ReadOnlyFile struct {
	name       string
	path       string
	size       int64
	createdAt  time.Time
	modifiedAt time.Time
	data       []byte
}

func (r *ReadOnlyFile) Name() string          { return r.name }
func (r *ReadOnlyFile) Path() string          { return r.path }
func (r *ReadOnlyFile) Size() int64           { return int64(len(r.data)) }
func (r *ReadOnlyFile) CreatedAt() time.Time  { return r.createdAt }
func (r *ReadOnlyFile) ModifiedAt() time.Time { return r.modifiedAt }
func (r *ReadOnlyFile) Read(p []byte) (int, error) {
	copy(p, r.data)
	return len(r.data), nil
}

type VirtualFileSystem struct {
	root *Katalog
}

func NewVirtualFileSystem() *VirtualFileSystem {
	return &VirtualFileSystem{
		root: &Katalog{name: "root", path: "/", items: make(map[string]FileSystemItem)},
	}
}

func (vfs *VirtualFileSystem) CreateFile(path, name string, data []byte) error {
	folder, err := vfs.FindFolder(path)
	if err != nil {
		return err
	}
	return folder.AddItem(&Plik{name: name, path: path + name, data: data, createdAt: time.Now(), modifiedAt: time.Now()})
}

func (vfs *VirtualFileSystem) CreateFolder(path, name string) error {
	folder, err := vfs.FindFolder(path)
	if err != nil {
		return err
	}
	return folder.AddItem(&Katalog{name: name, path: path + name + "/", items: make(map[string]FileSystemItem), createdAt: time.Now(), modifiedAt: time.Now()})
}

func (vfs *VirtualFileSystem) FindItem(path string) (FileSystemItem, error) {
	return vfs.findItemRecursive(vfs.root, path)
}

func (vfs *VirtualFileSystem) FindFolder(path string) (*Katalog, error) {
	item, err := vfs.FindItem(path)
	if err != nil {
		return nil, err
	}
	folder, ok := item.(*Katalog)
	if !ok {
		return nil, ErrNotDirectory
	}
	return folder, nil
}

func (vfs *VirtualFileSystem) findItemRecursive(folder *Katalog, path string) (FileSystemItem, error) {
	if folder.Path() == path {
		return folder, nil
	}
	for _, item := range folder.items {
		if item.Path() == path {
			return item, nil
		}
		if subFolder, ok := item.(*Katalog); ok {
			if found, err := vfs.findItemRecursive(subFolder, path); err == nil {
				return found, nil
			}
		}
	}
	return nil, ErrItemNotFound
}

func (vfs *VirtualFileSystem) DeleteItem(path string) error {
	parentPath, name := splitPath(path)
	folder, err := vfs.FindFolder(parentPath)
	if err != nil {
		return err
	}
	return folder.RemoveItem(name)
}

func splitPath(path string) (string, string) {
	if path == "/" {
		return "/", ""
	}
	lastSlash := len(path) - 1
	for lastSlash >= 0 && path[lastSlash] != '/' {
		lastSlash--
	}
	if lastSlash == 0 {
		return "/", path[1:]
	}
	return path[:lastSlash+1], path[lastSlash+1:]
}

func (vfs *VirtualFileSystem) CreateSymlink(path, name, pathOriginal string) error {
	folder, err := vfs.FindFolder(path)
	if err != nil {
		return err
	}
	original, err := vfs.FindItem(pathOriginal)
	if err != nil {
		return err
	}
	return folder.AddItem(&SymLink{name: name, path: path + name, createdAt: time.Now(), modifiedAt: time.Now(), target: original})
}

func main() {
	vfs := NewVirtualFileSystem()
	vfs.CreateFolder("/", "docs")
	vfs.CreateFile("/docs/", "file1.txt", []byte("Hello, world!"))
	vfs.CreateSymlink("/", "linkToDocs", "/docs/file1.txt")

	item, err := vfs.FindItem("/docs/file1.txt")
	if err == nil {
		fmt.Println("Found item:", item.Name())
	} else {
		fmt.Println("Error finding item:", err)
	}

	vfs.DeleteItem("/docs/file1.txt")
	item, err = vfs.FindItem("/docs/file1.txt")
	if err == nil {
		fmt.Println("Found item:", item.Name())
	} else {
		fmt.Println("Error finding item:", err)
	}

}
