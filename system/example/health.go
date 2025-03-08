package main

import (
	"fmt"


	"architect/system"
)

func main() {
	mem := system.ReadMemoryStats()
	disk, err := system.ReadDiskInfo("./")
	if err != nil {
		panic(err)
	}

	fmt.Println("In MegaBytes:")
	fmt.Printf("MemTotal: %d | MemFree: %d | MemAvailable: %d\n", mem.TotalMem/1024/1024, mem.FreeMem/1024, mem.AvailableMem/1024/1024)
	fmt.Println("In MegaBytes:")
	fmt.Printf("DiskTotal: %d | DiskFree: %d | DiskAvailable: %d\n", disk.TotalDisk/1024/1024, disk.FreeDisk/1024/1024, disk.AvailableDisk/1024/1024)
}