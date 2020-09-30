package util

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/devops-salt/src/message"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var (
 timeout = 30 * time.Second
 task_timeout int64 = 30
 tmp_dir = "/tmp/"
)

// FileExist determine if the file exists.
func FileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

// WritePidFile write pid into the file.
func WritePidFile(filename string) error {
	dir, _ := filepath.Split(filename)
	if !FileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(filename, []byte(strconv.Itoa(os.Getpid())), 0666); err != nil {
		return err
	}

	return nil
}

// GetLocalIP
func GetLocalIP() (ipv4s []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ipv4s = append(ipv4s, ipnet.IP.String())
			}
		}
	}
	return ipv4s
}

func WriteShellFile(filename string, cmd string)  error {
	dir, _ := filepath.Split(filename)
	if !FileExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	if err := ioutil.WriteFile(filename, []byte(cmd), 0666); err != nil {
		return err
	}

	return nil

}

func DelFile(filename string) error {
	if err := os.Remove(filename); err != nil {
		return err
	}
	return nil
}

// Compatible with complex commands by executing shell
func command(issue string, task_id string)  (string, error){
	filename := fmt.Sprintf("%v%v.sh",tmp_dir, task_id)
	if err := WriteShellFile(filename, issue); err != nil {
		return "", err
	}
	cmd := exec.Command("sh", filename)
	defer DelFile(filename)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	time.AfterFunc(timeout,
		func() {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		})

	if err := cmd.Start(); err != nil {
		return string(out.Bytes()), err
	}

	if err := cmd.Wait(); err != nil {
		return string(out.Bytes()), err
	}
	return string(out.Bytes()), nil

}

func Command(msg *message.Task, callback func(out *message.Callback, host string)){
	res := &message.Callback{
		Task_id: msg.Task_id,
		Code: 500,
	}
	iss, err := Base64decode(msg.Issue)
	if err == nil {
		out, err := command(iss, msg.Task_id)
		if err != nil {
			res.Data = err.Error()
		} else {
			res.Data = out
			res.Code = 200
		}
	} else {
		res.Data = err.Error()
	}
	if callback != nil {
		callback(res, msg.Callback)
	}
}

//func command(cmd string) (string, error){
//	var Timeout = 10 * time.Second
//	ctxt, cancel := context.WithTimeout(context.Background(), Timeout)
//	defer cancel()
//	issue := strings.Split(cmd, " ")
//	c := exec.Command(issue[0], issue[1:]...)
//	if out, err := c.Output(); err != nil {
//		// 检查错误原因
//		if ctxt.Err() != nil && ctxt.Err() == context.DeadlineExceeded {
//			return "nil", errors.New("command timeout")
//		}
//		return "", err
//	} else {
//		return string(out), nil
//	}
//}



func Base64decode(base string) (string, error)  {
	deco, err := base64.StdEncoding.DecodeString(base)
	if (err != nil) {
		return  "", err
	}
	return string(deco), nil
}

func Bsse64encode(base string)  (string){
	return base64.StdEncoding.EncodeToString([]byte(base))

}

// Check task
func CheckTask(p *message.Package, source []string) bool {
	dst_addr := false
	for _, dst := range p.Source {
		for _, rsc := range source {
			if dst == rsc {
				dst_addr = true
			}
		}
	}
	// The first address is the NAT address
	task_id := fmt.Sprintf("%v%v%d",p.Source[0], *p.Issue, *p.Timestamp)
	return len(p.Source) >= 2 && dst_addr && Md5(task_id) == *p.TaskId && time.Now().Unix() - *p.Timestamp < task_timeout
}

func CheckHeartbeat(p *message.Package) bool {
	return p.Attributes.Number() == message.Package_HEARTBEAT.Number() && time.Now().Unix() - *p.Timestamp < task_timeout
}

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

