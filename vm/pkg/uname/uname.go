package uname

import (
	"syscall"
)

func charsToString(ca []int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}

// Uname is wrapper for syscall uname
type Uname struct {
	ub syscall.Utsname
}

// Init call uname syscall
func New() (*Uname, error) {
	u := &Uname{}
	err := syscall.Uname(&u.ub)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// Machine return machine hardware name (like uname -m)
func (u *Uname) Machine() string {
	return charsToString(u.ub.Machine[:])
}

// Sysname return operating system (like uname -s)
func (u *Uname) Sysname() string {
	return charsToString(u.ub.Sysname[:])
}

// Nodename return network node hostname (like uname -n)
func (u *Uname) Nodename() string {
	return charsToString(u.ub.Nodename[:])
}

// KernelRelease return kernel release (like uname -r)
func (u *Uname) KernelRelease() string {
	return charsToString(u.ub.Release[:])
}

// KernelVersion return kernel version (like uname -v)
func (u *Uname) KernelVersion() string {
	return charsToString(u.ub.Version[:])
}
