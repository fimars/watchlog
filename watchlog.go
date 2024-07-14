package watchdog

import (
	"fmt"
	"os"
	"slices"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// 列表
type Logs struct {
	entry []uint64
}

func (l *Logs) Add(name uint64) {
	l.entry = append(l.entry, name)
}
func (l *Logs) Exist(name uint64) bool {
	return slices.Contains(l.entry, name)
}

// 监控者
type Watcher interface {
	Watch(name string, op fsnotify.Op)
}
type WatchDogs struct {
	DebugLineReader *func(line string)
	logs            *Logs
}

var _ Watcher = (*WatchDogs)(nil)

func NewWatch() *WatchDogs {
	return &WatchDogs{
		logs: &Logs{},
	}
}

func (w *WatchDogs) Watch(name string, op fsnotify.Op) {
	inode, err := getInode(name)
	if err != nil {
		fmt.Println("failed to get inode information")
	}

	if w.logs.Exist(inode) {
		return
	}
	w.logs.Add(inode)

	// printTime("%v", w.logs.entry)

	readFunc := ReadErrorLine
	if w.DebugLineReader != nil {
		readFunc = *w.DebugLineReader
	}

	go Tail(name, readFunc, op)
}

func getInode(filePath string) (uint64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("failed to get inode information")
	}

	return stat.Ino, nil
}
