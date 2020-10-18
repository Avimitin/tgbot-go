package hardwareInfo

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"time"
)

func GetCpuModel() (string, error) {
	cpuInfo, err := cpu.Info()

	if err != nil {
		return "Fail loading.", err
	}

	var text string

	for _, ci := range cpuInfo {
		text += ci.ModelName
	}

	return text, nil
}

func GetCpuPercent() (string, error) {
	percents, err := cpu.Percent(time.Second, true)

	if err != nil { return "Fail loading.", err }

	var eachCPUInfo string
	for id, percent := range percents {
		eachCPUInfo += fmt.Sprintf("CPU %d: %.2f\n", id, percent)
	}

	return fmt.Sprintf("CPU Percentage: \n%s", eachCPUInfo), nil
}

func GetCpuLoad() (string, error) {
	info, err := load.Avg()
	if err != nil { return "Fail loading.", err }
	return fmt.Sprintf("CPU Load: %v", info), err
}

func GetDiskUsage(path string) (string, error) {
	usageStat, err := disk.Usage(path)
	if err != nil { return "Fail loading.", err }

	return fmt.Sprintf(
		"DiskPath: %v\nTotal: %v GB\nFree: %v GB\nUsed: %v GB\nPercent: %.2f %%",
		usageStat.Path, usageStat.Total/1000000000, usageStat.Free/1000000000, usageStat.Used/1000000000, usageStat.UsedPercent),
		nil
}

func GetMemUsage() (string, error) {
	swapMem, err := mem.VirtualMemory()
	if err != nil { return "Fail loading.", err }

	return fmt.Sprintf("Total: %v GB\nFree: %v GB\nUsed: %v GB\nPercent: %.2f %%",
		swapMem.Total/1000000000, swapMem.Free/1000000000, swapMem.Used/1000000000, swapMem.UsedPercent),
	nil
}