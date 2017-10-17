package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	flag "github.com/spf13/pflag"
)

type selpg_args struct {
	start_page          int
	end_page            int
	in_filename         string
	page_len            int
	form_feed_delimited bool
	print_dest          string
}

var progname string
var err error

func main() {
	progname = os.Args[0]

	var sa selpg_args
	processArgs(&sa)
	// fmt.Printf("sa = %+v\n", sa)
	processInput(&sa)
}

func errExit() {
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
		os.Exit(1)
	}
}

func processArgs(sa *selpg_args) {
	flag.IntVarP(&sa.start_page, "start-page", "s", 1, "start page number")
	flag.IntVarP(&sa.end_page, "end-page", "e", 1, "end page number")
	flag.IntVarP(&sa.page_len, "page-length", "l", 72, "lines per page")
	flag.BoolVarP(&sa.form_feed_delimited, "form-feed-delimited", "f", false, "form feed delimited")
	flag.StringVarP(&sa.print_dest, "print-dest", "d", "", "print dest")
	flag.Parse()
	sa.in_filename = ""
	if len(flag.Args()) > 0 {
		sa.in_filename = flag.Args()[0]
	}
}

func processInput(sa *selpg_args) {
	var page_count int

	reader := makeReader(sa)
	writer, sub_proc := makeWriter(sa)

	if !sa.form_feed_delimited {
		page_count = pagingL(reader, &writer, sa)
	} else {
		page_count = pagingF(reader, &writer, sa)
	}

	warnAtPageCount(page_count, sa)

	if sub_proc != nil {
		// 关闭管道写端以使子进程终止
		writer.(io.WriteCloser).Close()
		// 回收子进程
		sub_proc.Wait()
	}
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}

func makeReader(sa *selpg_args) *bufio.Reader {
	in_fd := os.Stdin
	if len(sa.in_filename) > 0 {
		in_fd, err = os.Open(sa.in_filename)
		errExit()
	}
	return bufio.NewReaderSize(in_fd, 16*1024)
}

func makeWriter(sa *selpg_args) (io.Writer, *exec.Cmd) {
	var writer io.Writer
	var sub_proc *exec.Cmd

	writer = os.Stdout
	if len(sa.print_dest) > 0 {
		// sub_proc := exec.Command("lp", "-d", print_dest)
		// 使用"cat -n"代替
		sub_proc = exec.Command("cat", "-n")
		writer, err = sub_proc.StdinPipe()
		errExit()
		sub_proc.Stdout = os.Stdout
		sub_proc.Stderr = os.Stderr
		sub_proc.Start()
	}
	return writer, sub_proc
}

func pagingL(reader *bufio.Reader, writer *io.Writer, sa *selpg_args) int {
	var line string
	line_ctr := 0
	page_ctr := 1

	for {
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		errExit()

		line_ctr++
		if line_ctr > sa.page_len {
			page_ctr++
			line_ctr = 1
		}
		if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
			fmt.Fprintf(*writer, line)
		}
	}
	return page_ctr
}

func pagingF(reader *bufio.Reader, writer *io.Writer, sa *selpg_args) int {
	var page string
	page_ctr := 1

	for {
		page, err = reader.ReadString('\f')
		if err == io.EOF {
			break
		}
		errExit()

		page_ctr++
		if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
			fmt.Fprintf(*writer, page)
		}
	}
	return page_ctr
}

func warnAtPageCount(page_count int, sa *selpg_args) {
	if page_count < sa.start_page {
		fmt.Fprintf(os.Stderr, "%s: start_page (%d) greater than total pages (%d),"+
			" no output written\n",
			progname, sa.start_page, page_count)
	} else if page_count < sa.end_page {
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d),"+
			" less output than expected\n",
			progname, sa.end_page, page_count)
	}
}
