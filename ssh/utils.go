package ssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
	"k8s.io/klog/v2"
)

func fileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func readPipe(stdout, stderr io.Reader, isStdout bool) {
	var combineSlice []string
	var combineLock sync.Mutex
	doneout := make(chan error, 1)
	doneerr := make(chan error, 1)
	go func() {
		doneerr <- readpipe(stderr, &combineSlice, &combineLock, isStdout)
	}()
	go func() {
		doneout <- readpipe(stdout, &combineSlice, &combineLock, isStdout)
	}()
	<-doneerr
	<-doneout
}

func readpipe(pipe io.Reader, combineSlice *[]string, combineLock *sync.Mutex, isStdout bool) error {
	r := bufio.NewReader(pipe)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return err
		}

		combineLock.Lock()
		*combineSlice = append(*combineSlice, string(line))
		klog.Infof("Command execution result is: %s", line)
		if isStdout {
			fmt.Println(string(line))
		}
		combineLock.Unlock()
	}
}

func WaitSSHReady(ssh Interface, tryTimes int) error {
	var err error
	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		for i := 0; i < tryTimes; i++ {
			err = ssh.Ping()
			if err == nil {
				return nil
			}
			klog.Errorf("SSH connection failed (%v/%v), try again, error: %s", i+1, tryTimes, err)
			time.Sleep(time.Duration(1) * time.Second)
		}
		return fmt.Errorf("wait for [%v] ssh ready timeout: %v, ensure that the IP address or password is correct", ssh.GetIP(), err)
	})
	return eg.Wait()
}
