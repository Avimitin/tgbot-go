package hardwareInfo

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
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

	if err != nil {
		return "Fail loading.", err
	}

	var eachCPUInfo string
	for id, percent := range percents {
		eachCPUInfo += fmt.Sprintf("CPU %d: %.2f\n", id, percent)
	}

	return fmt.Sprintf("CPU Percentage: \n%s", eachCPUInfo), nil
}

func GetCpuLoad() (string, error) {
	info, err := load.Avg()
	return fmt.Sprintf("CPU Load: %v", info), err
}