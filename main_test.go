package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Tests POST request with correct request body
func TestDataHandlerSuccess(t *testing.T) {
	body := `
	{
		"level": "full",
		"minimal": {
			"current_os_version": "31.20191122.dev.0",
			"instance_type": null,
			"original_os_version": "",
			"platform": "qemu"
		},
		"full": {
			"container_rt": {
				"crio": {
					"is_running": false,
					"num_containers": 0
				},
				"docker": {
					"is_running": false,
					"num_containers": 0
				},
				"podman": {
					"is_running": true,
					"num_containers": 0
				},
				"systemd_nspawn": {
					"is_running": false,
					"num_containers": 0
				}
			},
			"hardware": {
				"cpu": {
					"lscpu": [
						{
							"data": "x86_64",
							"field": "Architecture:"
						},
						{
							"data": "32-bit, 64-bit",
							"field": "CPU op-mode(s):"
						},
						{
							"data": "Little Endian",
							"field": "Byte Order:"
						},
						{
							"data": "40 bits physical, 48 bits virtual",
							"field": "Address sizes:"
						},
						{
							"data": "8",
							"field": "CPU(s):"
						},
						{
							"data": "0-7",
							"field": "On-line CPU(s) list:"
						},
						{
							"data": "1",
							"field": "Thread(s) per core:"
						},
						{
							"data": "1",
							"field": "Core(s) per socket:"
						},
						{
							"data": "8",
							"field": "Socket(s):"
						},
						{
							"data": "1",
							"field": "NUMA node(s):"
						},
						{
							"data": "GenuineIntel",
							"field": "Vendor ID:"
						},
						{
							"data": "6",
							"field": "CPU family:"
						},
						{
							"data": "142",
							"field": "Model:"
						},
						{
							"data": "Intel(R) Core(TM) i7-8650U CPU @ 1.90GHz",
							"field": "Model name:"
						},
						{
							"data": "10",
							"field": "Stepping:"
						},
						{
							"data": "2111.998",
							"field": "CPU MHz:"
						},
						{
							"data": "4223.99",
							"field": "BogoMIPS:"
						},
						{
							"data": "VT-x",
							"field": "Virtualization:"
						},
						{
							"data": "KVM",
							"field": "Hypervisor vendor:"
						},
						{
							"data": "full",
							"field": "Virtualization type:"
						},
						{
							"data": "256 KiB",
							"field": "L1d cache:"
						},
						{
							"data": "256 KiB",
							"field": "L1i cache:"
						},
						{
							"data": "32 MiB",
							"field": "L2 cache:"
						},
						{
							"data": "128 MiB",
							"field": "L3 cache:"
						},
						{
							"data": "0-7",
							"field": "NUMA node0 CPU(s):"
						},
						{
							"data": "Mitigation; PTE Inversion; VMX conditional cache flushes, SMT disabled",
							"field": "Vulnerability L1tf:"
						},
						{
							"data": "Mitigation; Clear CPU buffers; SMT Host state unknown",
							"field": "Vulnerability Mds:"
						},
						{
							"data": "Mitigation; PTI",
							"field": "Vulnerability Meltdown:"
						},
						{
							"data": "Mitigation; Speculative Store Bypass disabled via prctl and seccomp",
							"field": "Vulnerability Spec store bypass:"
						},
						{
							"data": "Mitigation; usercopy/swapgs barriers and __user pointer sanitization",
							"field": "Vulnerability Spectre v1:"
						},
						{
							"data": "Mitigation; Full generic retpoline, IBPB conditional, IBRS_FW, STIBP disabled, RSB filling",
							"field": "Vulnerability Spectre v2:"
						},
						{
							"data": "fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon rep_good nopl xtopology cpuid tsc_known_freq pni pclmulqdq vmx ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch cpuid_fault invpcid_single pti ssbd ibrs ibpb tpr_shadow vnmi flexpriority ept vpid ept_ad fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx rdseed adx smap clflushopt xsaveopt xsavec xgetbv1 xsaves arat umip md_clear",
							"field": "Flags:"
						}
					]
				},
				"disk": {
					"blockdevices": [
						{
							"children": null,
							"fsavail": null,
							"fstype": null,
							"fsuse%": null,
							"label": null,
							"mountpoint": null,
							"name": "sr0"
						},
						{
							"children": [
								{
									"children": null,
									"fsavail": "264.8M",
									"fstype": "ext4",
									"fsuse%": "21%",
									"label": "boot",
									"mountpoint": "/boot",
									"name": "vda1"
								},
								{
									"children": null,
									"fsavail": "123.9M",
									"fstype": "vfat",
									"fsuse%": "2%",
									"label": "EFI-SYSTEM",
									"mountpoint": "/boot/efi",
									"name": "vda2"
								},
								{
									"children": null,
									"fsavail": null,
									"fstype": null,
									"fsuse%": null,
									"label": null,
									"mountpoint": null,
									"name": "vda3"
								},
								{
									"children": null,
									"fsavail": "5.8G",
									"fstype": "xfs",
									"fsuse%": "23%",
									"label": "root",
									"mountpoint": "/sysroot",
									"name": "vda4"
								}
							],
							"fsavail": null,
							"fstype": null,
							"fsuse%": null,
							"label": null,
							"mountpoint": null,
							"name": "vda"
						}
					]
				},
				"memory": {
					"memory": [
						{
							"block": "0",
							"removable": false,
							"size": "128M",
							"state": "online"
						},
						{
							"block": "1-5",
							"removable": true,
							"size": "640M",
							"state": "online"
						},
						{
							"block": "6",
							"removable": false,
							"size": "128M",
							"state": "online"
						},
						{
							"block": "7-9",
							"removable": true,
							"size": "384M",
							"state": "online"
						},
						{
							"block": "10",
							"removable": false,
							"size": "128M",
							"state": "online"
						},
						{
							"block": "11",
							"removable": true,
							"size": "128M",
							"state": "online"
						},
						{
							"block": "12-15",
							"removable": false,
							"size": "512M",
							"state": "online"
						}
					]
				}
			},
			"network": {
				"GENERAL.CON-PATH": "--",
				"GENERAL.CONNECTION": "--",
				"GENERAL.DEVICE": "lo",
				"GENERAL.HWADDR": "00:00:00:00:00:00",
				"GENERAL.MTU": "65536",
				"GENERAL.STATE": "10 (unmanaged)",
				"GENERAL.TYPE": "loopback",
				"IP4.ADDRESS[1]": "127.0.0.1/8",
				"IP4.DNS[1]": "10.0.2.3",
				"IP4.GATEWAY": "--",
				"IP4.ROUTE[1]": "dst = 0.0.0.0/0, nh = 10.0.2.2, mt = 100",
				"IP4.ROUTE[2]": "dst = 0.0.0.0/0, nh = 10.0.2.2, mt = 0",
				"IP4.ROUTE[3]": "dst = 10.0.2.0/24, nh = 0.0.0.0, mt = 0",
				"IP4.ROUTE[4]": "dst = 10.0.2.0/24, nh = 0.0.0.0, mt = 100",
				"IP6.ADDRESS[1]": "::1/128",
				"IP6.ADDRESS[2]": "fec0::5054:ff:fe12:3456/64",
				"IP6.ADDRESS[3]": "fe80::5054:ff:fe12:3456/64",
				"IP6.GATEWAY": "--",
				"IP6.ROUTE[1]": "dst = ::1/128, nh = ::, mt = 256",
				"IP6.ROUTE[2]": "dst = ::/0, nh = fe80::2, mt = 100",
				"IP6.ROUTE[3]": "dst = fe80::/64, nh = ::, mt = 100",
				"IP6.ROUTE[4]": "dst = fe80::/64, nh = ::, mt = 256",
				"IP6.ROUTE[5]": "dst = fec0::/64, nh = ::, mt = 256",
				"IP6.ROUTE[6]": "dst = ::/0, nh = fe80::2, mt = 1024",
				"IP6.ROUTE[7]": "dst = ff00::/8, nh = ::, mt = 256, table=255",
				"WIRED-PROPERTIES.CARRIER": "on"
			}
		}
	}`
	req, err := http.NewRequest("POST", "/v1", strings.NewReader(body))
	req.Header.Set("Content-type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(dataHandlerV1)

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		b, _ := ioutil.ReadAll(rr.Body)
		t.Errorf("response body: %v", string(b))
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}
