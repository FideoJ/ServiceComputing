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
	page_len    int /* default value, can be overriden by "-l number" on command
	   line */
	page_type int /* 'l' for lines-delimited, 'f' for form-feed-delimited */
	/* default is 'l' */
	print_dest string
}

var progname string

func main() {
	progname = os.Args[0]

	var sa selpg_args
	process_args(&sa)

	fmt.Printf("sa = %+v\n", sa)

	process_input(&sa)
}

func err_exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL ERROR: %v\n", err)
		os.Exit(1)
	}
}

func process_args(sa *selpg_args) {
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

func process_input(sa *selpg_args) {
	var in_fd *os.File
	var reader *bufio.Reader
	var writer io.WriteCloser
	var output io.ReadCloser
	var err error
	var line_ctr, page_ctr int

	if sa.in_filename == "" {
		in_fd = os.Stdin
	} else {
		in_fd, err = os.Open(sa.in_filename)
		err_exit(err)
		reader = bufio.NewReaderSize(in_fd, 16*1024)
	}

	if sa.print_dest == "" {
		writer = os.Stdout
	} else {
		// cmd := exec.Command("lp -d" + sa.print_dest)
		cmd := exec.Command("wc", "-l")
		writer, err = cmd.StdinPipe()
		err_exit(err)
		output, err = cmd.StdoutPipe()
		cmd.Start()
	}

	if sa.page_type == 'l' {
		line_ctr = 0
		page_ctr = 1

		for {
			line, is_prefix, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			err_exit(err)

			if !is_prefix {
				line_ctr++
			}
			if line_ctr > sa.page_len {
				page_ctr++
				line_ctr = 1
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				writer.Write(line)
			}
		}
	} else {
		page_ctr = 1
		for {
			c, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			err_exit(err)

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

	if sa.print_dest != "" {
		var n int
		out := make([]byte, 16*1024)
		n, err = output.Read(out)
		err_exit(err)

		fmt.Println("OUTPUT IN DEST:", string(out[:n]))
	}
	fmt.Fprintf(os.Stderr, "%s: done\n", progname)
}
