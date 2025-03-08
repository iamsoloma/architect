package system

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

type Memory struct {
	TotalMem     uint64
	FreeMem      uint64
	AvailableMem uint64
}

type Disk struct {
	ID            string
	TotalDisk     uint64
	FreeDisk      uint64
	AvailableDisk uint64
	DiskType      string
}

func ReadMemoryStats() Memory {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	bufio.NewScanner(file)
	scanner := bufio.NewScanner(file)
	res := Memory{}
	for scanner.Scan() {
		key, value := parseLine(scanner.Text())
		switch key {
		case "MemTotal":
			res.TotalMem = value*1024
		case "MemFree":
			res.FreeMem = value*1024
		case "MemAvailable":
			res.AvailableMem = value*1024
		}
	}
	return res
}

func ReadDiskInfo(path string) (Disk, error) {
	var info unix.Statfs_t

	err := unix.Statfs(path, &info)
	if err != nil {
		return Disk{}, err
	}

	res := Disk{}
	res.AvailableDisk = info.Bavail * uint64(info.Bsize)
	res.TotalDisk = info.Blocks * uint64(info.Bsize)
	res.FreeDisk = info.Bfree * uint64(info.Bsize)

	return res, nil
}

func parseLine(raw string) (key string, value uint64) {
	text := strings.ReplaceAll(raw[:len(raw)-2], " ", "")
	keyValue := strings.Split(text, ":")
	if keyValue[1] == "" {
		return keyValue[0], 0
	}
	strconv.ParseUint(keyValue[1], 10, 64)

	return keyValue[0], toUint(keyValue[1])
}

func toUint(raw string) uint64 {
	if raw == "" {
		return 0
	}
	res, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}
