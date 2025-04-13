package main

// 原生epoll使用示例（Linux系统）
/*
#include <sys/epoll.h>
*/
import "C"

type EpollWrapper struct {
	epfd int
}

func NewEpoll() (*EpollWrapper, error) {
	epfd, err := C.epoll_create1(0)
	if err != nil {
		return nil, err
	}
	return &EpollWrapper{epfd: int(epfd)}, nil
}

func (e *EpollWrapper) Add(fd int) error {
	event := C.struct_epoll_event{
		events: C.EPOLLIN,
		data:   C.union_epoll_data{fd: C.int(fd)},
	}
	_, err := C.epoll_ctl(C.int(e.epfd), C.EPOLL_CTL_ADD, C.int(fd), &event)
	return err
}
