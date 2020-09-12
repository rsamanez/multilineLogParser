// Copyright (c) 2015 HPE Software Inc. All rights reserved.
// Copyright (c) 2013 ActiveState Software Inc. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/hpcloud/tail"
	"github.com/marcsauter/single"
)

func args2config() (tail.Config, int64) {
	config := tail.Config{Follow: true}
	n := int64(0)
	maxlinesize := int(0)
	flag.Int64Var(&n, "n", 0, "tail from the last Nth location")
	flag.IntVar(&maxlinesize, "max", 0, "max line size")
	flag.BoolVar(&config.Follow, "f", false, "wait for additional data to be appended to the file")
	flag.BoolVar(&config.ReOpen, "F", false, "follow, and track file rename/rotation")
	flag.BoolVar(&config.Poll, "p", false, "use polling, instead of inotify")
	flag.Parse()
	if config.ReOpen {
		config.Follow = true
	}
	config.MaxLineSize = maxlinesize
	return config, n
}

func main() {
	config, n := args2config()
	if flag.NFlag() < 1 {
		fmt.Println("need one or more files as arguments")
		os.Exit(1)
	}

	if n != 0 {
		config.Location = &tail.SeekInfo{-n, os.SEEK_END}
	}
	s := single.New("multilinesLogParser")
	if err := s.CheckLock(); err != nil && err == single.ErrAlreadyRunning {
		log.Fatal("another instance of the app is already running, exiting")
	} else if err != nil {
		// Another error occurred, might be worth handling it as well
		log.Fatalf("failed to acquire exclusive app lock: %v", err)
	}
	defer s.TryUnlock()
	done := make(chan bool)
	for _, filename := range flag.Args() {
		go tailFile(filename, config, done)
	}

	for _, _ = range flag.Args() {
		<-done
	}
}

func check(s string)(string){
	if s != "" {
		return ","
	}else{
		return ""
	}
}

func getEntreComillas(s string)(string){
	split := strings.Split(s, "'")
	if len(split) == 3 {
		return split[1]
	}else{
		return ""
	}
}

func getNumber(s string)(string){
	split := strings.Split(s, "=")
	split2 := strings.Split(strings.TrimSpace(split[1]), " ")
	return split2[0]
}

func transforma(data []string) (string) {
	var salida string = ""
	for x := range data{
		if strings.Index(data[x],"Error") != -1{
			salida += check(salida) + "\"logType\":\"" + strings.TrimSpace(data[x])+"\""
		}
		if strings.Index(data[x],":RequestURL =") != -1{
			salida += check(salida) + "\"requestURL\":\"" + getEntreComillas(data[x]) +"\""
		}
		if strings.Index(data[x],":RequestURI  =") != -1{
			salida += check(salida) + "\"requestURI\":\"" + getEntreComillas(data[x]) +"\""
		}
		if strings.Index(data[x],":HTTPVersion =") != -1{
			salida += check(salida) + "\"httpVersion\":\"" + getEntreComillas(data[x]) +"\""
		}
		if strings.Index(data[x],":Method      =") != -1{
			salida += check(salida) + "\"method\":\"" + getEntreComillas(data[x]) +"\""
		}
		if strings.Index(data[x],"X-Original-HTTP-Status-Line") != -1{
			salida += check(salida) + "\"statusLine\":\"" + getEntreComillas(data[x]) +"\""
		}
		if strings.Index(data[x],"X-Original-HTTP-Status-Code") != -1{
			salida += check(salida) + "\"statusCode\":" + getNumber(data[x]) 
		}
		if strings.Index(data[x],":HTTPResponseHeader =") != -1{
			salida += check(salida) + "\"httpResponseHeader\":\"" + getEntreComillas(data[x]) +"\""
		}

	}
	if salida != "" {
		salida = "{" + salida + "}"
	}
	return(salida)
}

func tailFile(filename string, config tail.Config, done chan bool) {
	defer func() { done <- true }()
	t, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	marca := "=============================="
	var iniciando int = 0
	var buffer = []string{""}
	buffer = nil
	for line := range t.Lines {
		if strings.TrimSpace(line.Text) == marca {
			if iniciando == 0 {
				//fmt.Println("INICIO")
				iniciando = 1
			}else{
				//fmt.Println("FIN")
				iniciando = 0
				if buffer == nil{
					fmt.Println("INICIO")
					iniciando = 1
				}
				fmt.Println(transforma(buffer))
				buffer = nil
			}
		}
		if iniciando == 1 {
			if line.Text != marca {
				buffer = append(buffer, line.Text)
			}
		}
	}
	err = t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}
