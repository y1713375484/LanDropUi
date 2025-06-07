package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

// App struct
type App struct {
	ctx context.Context
}

type progressReader struct {
	reader    io.Reader
	total     int64           // 总大小（可选，用于显示百分比）
	readEn    int64           //当前读到的文件字节
	lastPrint int64           //上次读到的文件字节
	ctx       context.Context //上下文
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Listen(receiveCount int, receivePath string, ipAddress string) string {
	if receivePath != "" {
		_, err := os.Stat(receivePath)
		//检测要保存的路径是否存在
		if os.IsNotExist(err) {
			err := os.Mkdir(receivePath, 0777)
			if err != nil {
				fmt.Println(err)
				return err.Error()
			}
		}
	}

	// 监听端口
	listen, err := net.Listen("tcp", ipAddress)

	if err != nil {
		fmt.Println("Error listening:", err)
		return "当前端口监听失败，请检查端口是否被占用" + err.Error()
	}
	defer listen.Close()
	receiveChan := make(chan struct{}, receiveCount) //同时接收的文件数量
	for {
		// 接受客户端连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, a.ctx, receiveChan, receivePath) // 使用 goroutine 处理连接
	}
}

func handleConnection(conn net.Conn, ctx context.Context, receiveChan chan struct{}, receivePath string) {
	defer func() {
		conn.Close()
		<-receiveChan
	}()
	reader := bufio.NewReader(conn)

	// 读取文件名称
	fileName, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Println("Error reading file name:", err)
		return
	}
	fileName = fileName[:len(fileName)-1]

	//读取文件大小
	size, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("Error reading file size:", err)
			return
		}
		fmt.Println("Error reading file size:", err)
		return
	}
	size = size[:len(size)-1]
	sizeInt, _ := strconv.ParseInt(size, 10, 64)

	// 创建文件
	file, err := os.Create(filepath.Join(receivePath, fileName))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()
	//为当前连接传输的文件绑定一个uuid
	fileUUid := uuid.New()
	ctx = context.WithValue(ctx, "fileUUID", fileUUid.String())
	r := &progressReader{reader: reader, total: sizeInt, ctx: ctx}
	runtime.EventsEmit(r.ctx, "findFileName", map[string]string{
		"fileUUID": fileUUid.String(),
		"fileName": fileName,
	})
	receiveChan <- struct{}{}

	io.CopyN(file, r, sizeInt)

	//最终刷新进度条
	r.Finish()

}

// 实现reader接口
func (pw *progressReader) Read(p []byte) (int, error) {
	n, err := pw.reader.Read(p)
	pw.readEn += int64(n)

	// 统一阈值：如果知道总大小，按1%计算；否则按1MB固定间隔
	threshold := int64(0)
	if pw.total > 0 {
		threshold = pw.total / 100 // 1% of total
	} else {
		threshold = 1024 * 1024 // 1MB
	}

	if pw.readEn-pw.lastPrint > threshold || pw.lastPrint == 0 {
		percent := float64(pw.readEn) / float64(pw.total) * 100
		fmt.Printf("\rCopied: %d/%d bytes (%.f%%)", pw.readEn, pw.total, percent)
		runtime.EventsEmit(pw.ctx, "percent", map[string]interface{}{
			"fileUUID": pw.ctx.Value("fileUUID"),
			"percent":  percent,
		})
		pw.lastPrint = pw.readEn
	}

	return n, err
}

// 最终刷新进度
func (pw *progressReader) Finish() {
	percent := float64(pw.readEn) / float64(pw.total) * 100
	runtime.EventsEmit(pw.ctx, "percent", map[string]interface{}{
		"fileUUID": pw.ctx.Value("fileUUID"),
		"percent":  percent,
	})
	fmt.Printf("\rCopied: %d/%d bytes (%.f%%)", pw.readEn, pw.total, percent)
}

// 弹窗
func (a *App) Msgalert(msg string) {
	dialogOptions := runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "提示",
		Message: msg,
	}
	_, err := runtime.MessageDialog(a.ctx, dialogOptions)
	if err != nil {
		fmt.Println(err)
	}

}

func (a *App) ChooseReceivePath() string {
	dialogOptions := runtime.OpenDialogOptions{
		Title: "选择要保存文件的路径",
	}
	path, err := runtime.OpenDirectoryDialog(a.ctx, dialogOptions)
	if err != nil {
		fmt.Println(err)
	}
	return path
}
