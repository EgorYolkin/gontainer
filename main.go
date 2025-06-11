package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// main - точка входа в программу.
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Using: gontainer run <команда> [аргументы...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		fmt.Println("Unknown command. Use 'run'.")
		os.Exit(1)
	}
}

// parent set up and run child process in the new namespaces
func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	// SysProcAttr set up for isolation
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: "rootfs",
		Setsid: true,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("parent error:", err)
		os.Exit(1)
	}
}

// child into isolation environment
func child() {
	fmt.Printf("Запускаю дочерний процесс с PID: %d\n", os.Getpid())

	// Mount rootfs as bind mount
	if err := unix.Mount("rootfs", "rootfs", "", int(unix.MS_BIND|unix.MS_REC), unsafe.Pointer(nil)); err != nil {
		fmt.Println("ОШИБКА при монтировании rootfs:", err)
		os.Exit(1)
	}

	if err := os.MkdirAll("rootfs/oldrootfs", 0700); err != nil {
		fmt.Println("oldrootfs creating err:", err)
		os.Exit(1)
	}

	if err := unix.PivotRoot("rootfs", "rootfs/oldrootfs"); err != nil {
		fmt.Println("pivot_root error:", err)
		os.Exit(1)
	}

	// Gonna new dir
	if err := os.Chdir("/"); err != nil {
		fmt.Println("ОШИБКА при смене директории:", err)
		os.Exit(1)
	}

	// Mount old filesystem
	if err := unix.Unmount("/oldrootfs", unix.MNT_DETACH); err != nil {
		fmt.Println("ОШИБКА при размонтировании oldrootfs:", err)
		os.Exit(1)
	}

	// Delete old core dir
	if err := os.RemoveAll("/oldrootfs"); err != nil {
		fmt.Println("oldrootfs deletion error:", err)
		os.Exit(1)
	}

	// Run user command
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("child error:", err)
		os.Exit(1)
	}
}
