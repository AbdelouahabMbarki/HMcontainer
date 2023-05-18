package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what??")
	}
}

func run() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running the /proc/self/exe command:", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	// Set hostname of the new UTS namespace
	if err := syscall.Sethostname([]byte("HMcontainer")); err != nil {
		fmt.Println("Error setting hostname:", err)
		os.Exit(1)
	}

	// Change to the new root file system.
	if err := syscall.Chroot("/home/bloguser/homeMadeContainer/rootfs"); err != nil {
		fmt.Println("Error changing root:", err)
		os.Exit(1)
	}

	// Change working directory after changing the root.
	if err := os.Chdir("/"); err != nil {
		fmt.Println("Error changing working directory:", err)
		os.Exit(1)
	}

	// Mount proc. This needs to be done after chroot and chdir.
	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		fmt.Println("Error mounting proc:", err)
		os.Exit(1)
	}
	// Look up the user and group. Here we'll use "guest" as an example.
	userName := "guest"
	u, err := user.Lookup(userName)
	if err != nil {
		fmt.Println("Error looking up user:", err)
		os.Exit(1)
	}

	// Parse the found UID and GID.
	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		fmt.Println("Error parsing UID:", err)
		os.Exit(1)
	}
	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		fmt.Println("Error parsing GID:", err)
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set the UID and GID for the command.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid:         uint32(uid),
			Gid:         uint32(gid),
			NoSetGroups: true,
		},
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running the child command:", err)
		os.Exit(1)
	}
}
