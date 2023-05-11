package tuntap

import (
	"fmt"
	"net"
	"os/exec"
	"sync"
	"syscall"

	"path/filepath"

	"time"

	"github.com/gamexg/gotool/sysinfo"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const NETWORK_KEY = `SYSTEM\CurrentControlSet\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}`
const ADAPTER_KEY = `SYSTEM\CurrentControlSet\Control\Class\{4D36E972-E325-11CE-BFC1-08002BE10318}`
const DEVICE_COMPONENT_ID = "tap0901"

const WINXP_32_INSTALL_BAT = "winxp32/install.bat"
const WINXP_64_INSTALL_BAT = "winxp64/install.bat"
const WIN7_32_INSTALL_BAT = "win732/install.bat"
const WIN7_64_INSTALL_BAT = "win764/install.bat"

const WINXP_32_UNINSTALL_BAT = "winxp32/uninstall.bat"
const WINXP_64_UNINSTALL_BAT = "winxp64/uninstall.bat"
const WIN7_32_UNINSTALL_BAT = "win732/uninstall.bat"
const WIN7_64_UNINSTALL_BAT = "win764/uninstall.bat"

const TCPIP_LINKAGE = `SYSTEM\CurrentControlSet\Services\Tcpip\Linkage`

func TAP_CONTROL_CODE(request, method uint32) uint32 {
	return (34 << 16) | (0 << 14) | (request << 2) | method
}

var TAP_IOCTL_SET_MEDIA_STATUS = TAP_CONTROL_CODE(6, 0)
var TAP_IOCTL_CONFIG_TUN = TAP_CONTROL_CODE(10, 0)

type tun struct {
	handle           windows.Handle
	netCfgInstanceId string // {36961F56-38A5-49CB-8409-B0C9C0D1119D}
	name             string //以太网 3
	overlappedWx     windows.Overlapped
	rm               sync.Mutex
	overlappedRx     windows.Overlapped
	wm               sync.Mutex
}

func NewTun(addr net.IP, network net.IP, mask net.IP) (Tun, error) {

	rHevent, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("CreateEvent:%v", err)
	}

	wHevent, err := windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("CreateEvent:%v", err)
	}

	componentId := DEVICE_COMPONENT_ID
	netCfgInstanceId, err := getTuntapComponentId(componentId)
	if err != nil {

		err = f.newTunDevice(componentId)
		if err != nil {
			return nil, err
		}
		for i := 0; i < 10; i++ {
			netCfgInstanceId, err = getTuntapComponentId(componentId)
			if err == nil {
				break
			}
			time.Sleep(1 * time.Second)
		}
	}

	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get the ID card %v, tun. Run with administrator privileges to create the Tun network card。", err)
	}

	name := ""
	for i := 0; i < 10; i++ {
		name, err = getTunTapName(netCfgInstanceId)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf(`\\.\Global\%s.tap`, netCfgInstanceId)
	dpath := windows.StringToUTF16(path)

	tap, err := windows.CreateFile(
		&dpath[0],
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_SYSTEM|windows.FILE_FLAG_OVERLAPPED,
		0)

	if err != nil {
		return nil, err
	}

	var returnLen uint32
	var configTunParam []byte = append(addr.To4(), network.To4()...)
	configTunParam = append(configTunParam, mask.To4()...)

	// 配置ip后是 TAP 设备。
	if err = windows.DeviceIoControl(
		tap,
		TAP_IOCTL_CONFIG_TUN,
		&configTunParam[0],
		uint32(len(configTunParam)),
		&configTunParam[0],
		uint32(len(configTunParam)),
		&returnLen,
		nil); err != nil {
		return nil, err
	}

	err = setTunTapIP(name, addr, mask)
	if err != nil {
		return nil, err
	}

	inBuffer := []byte{0x01, 0x00, 0x00, 0x00}
	if err = windows.DeviceIoControl(
		tap,
		TAP_IOCTL_SET_MEDIA_STATUS,
		&inBuffer[0],
		uint32(len(inBuffer)),
		&inBuffer[0],
		uint32(len(inBuffer)),
		&returnLen,
		nil); err != nil {
		return nil, err
	}
	/*
		device, err := winio.MakeOpenFile(tap)
		if err != nil {
			return nil, fmt.Errorf("winio.makeOpenFile ,%v", err)
		}*/
	t := tun{
		handle: tap,
		//	device:           device,
		name:             name,
		netCfgInstanceId: netCfgInstanceId,
	}

	t.overlappedRx.HEvent = rHevent
	t.overlappedWx.HEvent = wHevent

	return &t, nil
}

func newTunDevice(componentId string) error {
	batPath := ""
	sysVer, err := sysinfo.GetSemanticVersion()
	if err != nil {
		return fmt.Errorf("Failed to create the Tun network card, unable to obtain the system version，%v", err)
	}
	is64, err := sysinfo.Is64Sys()
	if err != nil {
		return fmt.Errorf("Failed to create the Tun network card, unable to obtain the system version，%v", err)
	}

	if sysVer.Major >= 6 {
		// wind7
		if is64 {
			batPath = WIN7_64_INSTALL_BAT
		} else {
			batPath = WIN7_32_INSTALL_BAT
		}
	} else {
		if is64 {
			batPath = WINXP_64_INSTALL_BAT
		} else {
			batPath = WINXP_32_INSTALL_BAT
		}
	}
	batPath, err = filepath.Abs(batPath)
	if err != nil {
		return err
	}

	cmd := exec.Command("cmd.exe", "/C", batPath)
	cmd.Dir = filepath.Dir(batPath)
	o, err := cmd.Output()
	if err != nil {
		out, err := simplifiedchinese.GB18030.NewDecoder().Bytes(o)
		if err == nil {
			o = out
		}
		fmt.Println(string(o))
	}

	return nil
}

func getTuntapComponentId(componentId string) (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE,
		ADAPTER_KEY,
		registry.ENUMERATE_SUB_KEYS|registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	keyNames, err := k.ReadSubKeyNames(-1)
	if err != nil {
		return "", err
	}

	for _, keyName := range keyNames {
		n, _ := matchKey(k, keyName, componentId)
		if n != "" {
			return n, nil
		}
	}
	fmt.Println(keyNames)
	return "", fmt.Errorf("not found TAP device.")
}

func matchKey(zones registry.Key, keyName string, componentId string) (string, error) {
	k, err := registry.OpenKey(zones, keyName, registry.READ)
	if err != nil {
		return "", err
	}
	defer k.Close()

	cId, _, err := k.GetStringValue("ComponentId")
	if cId == componentId {
		netCfgInstanceId, _, err := k.GetStringValue("NetCfgInstanceId")
		if err != nil {
			return "", err
		}
		return netCfgInstanceId, nil
	}
	return "", fmt.Errorf("ComponentId != componentId")
}

func getTunTapName(netCfgInstanceId string) (string, error) {
	path := fmt.Sprintf("%s\\%s\\Connection", NETWORK_KEY, netCfgInstanceId)
	k, err := registry.OpenKey(windows.HKEY_LOCAL_MACHINE, path, registry.READ)
	if err != nil {
		return "", fmt.Errorf("OpenKey(%s):%s", path, err)
	}
	defer k.Close()

	name, _, err := k.GetStringValue("Name")
	if err != nil {
		return "", fmt.Errorf("(%s).GetStringValue(name):%s", path, err)
	}

	return name, nil
}

func setTunTapIP(name string, ip, mask net.IP) error {
	setip := exec.Command("netsh")
	setip.SysProcAttr = &syscall.SysProcAttr{}
	cmd := fmt.Sprintf(`netsh interface ip set address "%s" static %s %s`, name, ip.To4().String(), mask.To4().String())

	setip.SysProcAttr.CmdLine = cmd

	o, err := setip.Output()
	if err != nil {
		out, err := simplifiedchinese.GB18030.NewDecoder().Bytes(o)
		if err == nil {
			o = out
		}

		return fmt.Errorf("cmd:%s,err:%s,out:%s", cmd, err, string(o))
	}
	return nil
}

func (f *tunTapF) Remove() error {
	batPath := ""
	sysVer, err := sysinfo.GetSemanticVersion()
	if err != nil {
		return fmt.Errorf("创建 tun 网卡失败，无法获得系统版本，%v", err)
	}
	is64, err := sysinfo.Is64Sys()
	if err != nil {
		return fmt.Errorf("创建 tun 网卡失败，无法获得系统版本，%v", err)
	}

	if sysVer.Major >= 6 {
		// wind7
		if is64 {
			batPath = WIN7_64_UNINSTALL_BAT
		} else {
			batPath = WIN7_32_UNINSTALL_BAT
		}
	} else {
		if is64 {
			batPath = WINXP_64_UNINSTALL_BAT
		} else {
			batPath = WINXP_32_UNINSTALL_BAT
		}
	}
	batPath, err = filepath.Abs(batPath)
	if err != nil {
		return err
	}
	fmt.Println(batPath)
	cmd := exec.Command("cmd.exe", "/C", batPath)
	cmd.Dir = filepath.Dir(batPath)
	o, err := cmd.Output()
	if err != nil {
		out, err := simplifiedchinese.GB18030.NewDecoder().Bytes(o)
		if err == nil {
			o = out
		}
		fmt.Println(string(o))
	}
	return nil
}

func (t *tun) Read(buf []byte) (int, error) {
	var l uint32
	t.rm.Lock()
	defer t.rm.Unlock()
	t.overlappedRx.InternalHigh = 0

	windows.ReadFile(t.handle, buf, &l, &t.overlappedRx)
	if _, err := windows.WaitForSingleObject(t.overlappedRx.HEvent, windows.INFINITE); err != nil {
		return int(uint(t.overlappedRx.InternalHigh)), fmt.Errorf("Wait ReadFile:%v", err)
	}
	return int(uint(t.overlappedRx.InternalHigh)), nil
}

func (t *tun) Write(buf []byte) (int, error) {
	l := uint32(len(buf))
	// 如果移除锁，1W 次写入速度能由 314ms 降低到 265ms 。
	// 但是多线程会直接完蛋。
	t.wm.Lock()
	defer t.wm.Unlock()

	windows.WriteFile(t.handle, buf, &l, &t.overlappedWx)
	if _, err := windows.WaitForSingleObject(t.overlappedWx.HEvent, windows.INFINITE); err != nil {
		return 0, fmt.Errorf("Wait WriteFile:%v", err)
	}
	return int(l), nil
}

func (t *tun) GetName() string {
	return t.name
}

func (t *tun) Close() error {
	//	return t.device.Close()
	return windows.CloseHandle(t.handle)
}

func cmdRun(cmd string) error {
	setip := exec.Command("netsh")
	setip.SysProcAttr = &syscall.SysProcAttr{}

	setip.SysProcAttr.CmdLine = cmd

	o, err := setip.Output()
	if err != nil {
		out, err2 := simplifiedchinese.GB18030.NewDecoder().Bytes(o)
		if err2 == nil {
			o = out
		}

		return fmt.Errorf("cmd:%s,err:%s,out:%s", cmd, err, string(o))
	}
	return nil
}
