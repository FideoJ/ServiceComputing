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
	start_page  int
	end_page    int
	in_filename string
	page_len    int
	page_type   int
	print_dest  string
}

var progname string
var err error

func main() {
	progname = os.Args[0]

	var sa selpg_args
	processArgs(&sa)

	fmt.Printf("sa = %+v\n", sa)

	processInput(&sa)
}

func errExit() {
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
		os.Exit(1)
	}
}

func processArgs(sa *selpg_args) {
	flag.IntVarP(&sa.start_page, "start-page", "s", -1, "start page number")
	flag.IntVarP(&sa.end_page, "end-page", "e", -1, "end page number")
	//
	flag.IntVarP(&sa.page_len, "page-length", "l", 72, "lines per page")
	flag.IntVarP(&sa.page_type, "form-feed-delimited", "f", 'l', "form feed delimited")
	flag.StringVarP(&sa.print_dest, "print-dest", "d", "", "print dest")
	flag.Parse()
	sa.in_filename = ""
	if len(flag.Args()) > 0 {
		sa.in_filename = flag.Args()[0]
	}
}

func processInput(sa *selpg_args) {
	var writer io.Writer
	var sub_proc *exec.Cmd
	var line_ctr, page_ctr int

	in_fd := os.Stdin
	if len(sa.in_filename) > 0 {
		in_fd, err = os.Open(sa.in_filename)
		errExit()
	}
	reader := bufio.NewReaderSize(in_fd, 16*1024)

	writer = os.Stdout
	if len(sa.print_dest) > 0 {
		// sub_proc := exec.Command("lp", "-d", print_dest)
		sub_proc = exec.Command("wc", "-l")
		writer, err = sub_proc.StdinPipe()
		errExit()
		sub_proc.Stdout = os.Stdout
		sub_proc.Stderr = os.Stderr
		sub_proc.Start()
	}

	if sa.page_type == 'l' {
		line_ctr = 0
		page_ctr = 1

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			errExit()

			if line_ctr > sa.page_len {
				page_ctr++
				line_ctr = 1
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				fmt.Fprintf(writer, line)
			}
		}
	} else {
		page_ctr = 1
		for {
			c, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			errExit()

			if c == '\f' {
				page_ctr++
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				writer.Write([]byte{c})
			}
		}
	}

	if page_ctr < sa.start_page {
		fmt.Fprintf(os.Stderr, "%s: start_page (%d) greater than total pages (%d),"+
			" no output written\n",
			progname, sa.start_page, page_ctr)
	} else if page_ctr < sa.end_page {
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d),"+
			" less output than expected\n",
			progname, sa.end_page, page_ctr)
	}

	if sub_proc != nil {
		writer.(io.WriteCloser).Close()
		sub_proc.Wait()
	}
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}
