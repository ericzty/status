package main

import (
	"os/exec"
	"strings"
	"time"
)

func Monitoring(st status) {
	for {
		for _, t := range st.values {
			st.Mutex.Lock()
			switch t.Method {
			//https://github.com/go-ping/ping
			case "icmp":
				out4, _ := exec.Command("ping", "-4", "-c 1", "-W 1", t.Target).Output()
				t.Status4.Up = strings.Contains(string(out4), "64 bytes from")
				out6, _ := exec.Command("ping", "-6", "-c 1", "-W 1", t.Target).Output()
				t.Status6.Up = strings.Contains(string(out6), "64 bytes from")
			case "head":
				out4, _ := exec.Command("curl", "--connect-timeout", "2", "-4", "-I", "https://"+t.Target).Output()
				t.Status4.Up = strings.Contains(string(out4), "200")
				out6, _ := exec.Command("curl", "--connect-timeout", "2", "-6", "-I", "https://"+t.Target).Output()
				t.Status6.Up = strings.Contains(string(out6), "200")
			case "match":
				out4, _ := exec.Command("curl", "--connect-timeout", "2", "-4", "https://"+t.Target).Output()
				t.Status4.Up = strings.Contains(string(out4), t.Match)
				out6, _ := exec.Command("curl", "--connect-timeout", "2", "-6", "https://"+t.Target).Output()
				t.Status6.Up = strings.Contains(string(out6), t.Match)
			}
			go CheckStatus(t, &t.Status4)
			go CheckStatus(t, &t.Status6)
			st.Mutex.Unlock()
		}
		time.Sleep(5 * time.Second)
	}
}

func CheckStatus(t *target, ts *targetStatus) {
	if ts.Up == true {
		if ts.Fails != 0 {
			ts.Fails = 0
			go IsUp(t, ts)
		}
	} else {
		ts.Fails++
		if ts.Fails > 2 {
			go IsDown(t, ts)
		}
	}
}
