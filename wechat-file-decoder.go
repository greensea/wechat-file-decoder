package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	if len(os.Args) == 4 {
		Decode()
	} else if len(os.Args) == 2 {
		Guess()
	} else {
		Usage()
		os.Exit(0)
	}
}

func Usage() {
	fmt.Printf(`用法：
	* 扫描目录下的文件，猜测用于编解码的异或值：
		wechat-file-decoder <文件目录>
	* 给定一个异或值，对一个文件进行解码：
		wechat-fle-decoder <异或值> <微信文件路径> <解码后保存路径>
	`)
}

func Guess() {
	var fileCount int

	dirPath := os.Args[1]

	// 假定 jpeg 文件占大多数，我们尝试使用 jpeg 的文件头 0xffd8 对每一个文件的头两个字节进行异或操作
	// 若结果的两个字节均相同，我们就认为这可能是一个合法的 jpeg 文件，其值即为编解码使用的异或值
	type Code struct {
		B     byte
		Count int
	}
	stat := make(map[byte]*Code)
	filepath.Walk(dirPath, func(p string, info fs.FileInfo, err error) error {
		fi, err := os.Open(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "无法读取文件 %s: %v\n", p, err)
			return nil
		}

		fmt.Printf("正在扫描 %v\n", p)

		defer fi.Close()
		r := bufio.NewReader(fi)

		b1, _ := r.ReadByte()
		b2, _ := r.ReadByte()

		xorB1 := b1 ^ 0xff
		xorB2 := b2 ^ 0xd8

		if xorB1 != xorB2 {
			return nil
		}

		_, ok := stat[xorB1]
		if ok != true {
			stat[xorB1] = &Code{
				B:     xorB1,
				Count: 0,
			}
		}
		stat[xorB1].Count++
		fileCount++

		return nil
	})

	// 对结果进行排序
	var codes []*Code
	for _, v := range stat {
		codes = append(codes, v)
	}

	sort.Slice(codes, func(i, j int) bool {
		return codes[i].Count > codes[j].Count
	})

	fmt.Printf("\n====扫描结果====\n")
	fmt.Printf("可能的异或值 | 出现数量 | 占比\n")

	for _, v := range codes {
		fmt.Printf("%02x                %5d    %0.1f\n", v.B, v.Count, (float32(v.Count))/float32(fileCount)*100)
	}

	fmt.Printf("==============\n\n")
	if len(codes) == 0 {
		fmt.Printf("没有猜测到异或值，退出\n")
		os.Exit(0)
	}

	fmt.Printf("猜测出的异或值是: %2x\n", codes[0].B)
	fmt.Printf("请使用以下命令进行解码：\nwechat-file-decoder %02x <待解码文件> <输出文件> 进行解码\n", codes[0].B)

}

func Decode() {
	var xorB byte
	_, err := fmt.Sscanf(os.Args[1], "%02x", &xorB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "异或值错误: %v\n", err)
		return
	}

	from := os.Args[2]
	to := os.Args[3]

	fiR, err := os.Open(from)
	if err != nil {
		fmt.Fprintf(os.Stderr, "打开文件失败: %v\n", err)
		return
	}
	defer fiR.Close()

	r := bufio.NewReader(fiR)

	fiW, err := os.Create(to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "创建文件 %s 失败: %v\n", to, err)
		return
	}
	defer fiW.Close()

	w := bufio.NewWriter(fiW)

	buf := make([]byte, 128*1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Fprintf(os.Stderr, "读取文件失败: %v\n", err)
			break
		}
		if n == 0 {
			break
		}

		for i := 0; i < len(buf); i++ {
			buf[i] ^= xorB
		}

		w.Write(buf[:n])
	}
}
