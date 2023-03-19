package host

import "syscall"

func TotalMemory() uint64 {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		panic(err)
	}

	return uint64(info.Totalram)
}
