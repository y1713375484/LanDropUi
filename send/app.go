package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

// App struct
type App struct {
	ctx context.Context
}

type progressReader struct {
	reader    io.Reader
	total     int64 // 总大小（可选，用于显示百分比）
	readEn    int64
	lastPrint int64
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

// 选择文件
func (a *App) ChooseFile(lastFilePathList map[string]map[string]interface{}) map[string]map[string]interface{} {
	dialogOptions := runtime.OpenDialogOptions{
		Title: "选择要发送的文件",
	}
	filePaths, err := runtime.OpenMultipleFilesDialog(a.ctx, dialogOptions)
	if err != nil {
		fmt.Println("Error opening file dialog:", err)
	}

	filePathList := map[string]map[string]interface{}{}

	//判断上次文件传输列表是否为空
	if len(lastFilePathList) != 0 {
		filePathList = lastFilePathList
	}

	for _, filePath := range filePaths {
		fileUUid := uuid.New().String()
		filePathList[fileUUid] = map[string]interface{}{
			"filePath": filePath,
			"percent":  0,
		}
	}
	return filePathList
}

func (a *App) Send(filepathList map[string]map[string]interface{}, sendCount int, ipAddress string) string {
	sendChan := make(chan struct{}, sendCount)
	defer close(sendChan)
	sendErrChan := make(chan string, len(filepathList))
	waitGroup := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background()) // 创建可取消的context
	defer cancel()
	// 打开要传输的文件
	for fileUUid, filePathArray := range filepathList {
		//进度为0的文件才传输
		if filePathArray["percent"].(float64) == 0 {
			waitGroup.Add(1)
			go func(fp string) {
				defer waitGroup.Done()
				select {
				case sendChan <- struct{}{}:
					defer func() { <-sendChan }()
				case <-ctx.Done():
					return
				}

				// 检查是否已被取消
				select {
				case <-ctx.Done():
					return
				default:
				}

				if err := a.SendDo(fileUUid, fp, ipAddress); err != nil {
					cancel()
					sendErrChan <- err.Error()
				}

			}(filePathArray["filePath"].(string))

		}

	}

	errMsg := ""
	waitGroup.Wait()
	//关闭管道，否则读错误会阻塞
	close(sendErrChan)

	select {
	case errMsg = <-sendErrChan:
	default:
		errMsg = ""
	}
	if errMsg == "" {
		return "文件发送完毕"
	} else {
		return errMsg
	}
}

func (a *App) SendDo(fileUUid, filePath string, ipAddress string) error {
	// 为每个文件建立新的连接
	conn, err := net.Dial("tcp", ipAddress)
	if err != nil {
		return errors.New("请检查接收端是否以准备接收文件")
	}

	defer conn.Close()

	file, err := os.Open(filePath)
	if err != nil {

		return errors.New(file.Name() + "文件发送失败")

	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err

	}
	fileSize := fileInfo.Size()

	// 发送文件名称和大小
	fileName := filepath.Base(filePath)
	_, err = conn.Write([]byte(fileName + "\n" + strconv.FormatInt(fileSize, 10) + "\n"))
	if err != nil {
		return errors.New("文件信息传输失败")
	}
	//为当前连接传输的文件绑定一个uuid
	fileCtx := context.WithValue(a.ctx, "fileUUID", fileUUid)
	r := &progressReader{
		reader: file,
		total:  fileSize,
		ctx:    fileCtx,
	}

	_, err = io.Copy(conn, r)
	if err != nil {
		return errors.New("文件发送失败")
	}

	//最终刷新进度条
	r.Finish()
	return nil
}

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

func (pw *progressReader) Finish() {
	percent := float64(pw.readEn) / float64(pw.total) * 100
	fmt.Printf("\rCopied: %d/%d bytes (%.f%%)", pw.readEn, pw.total, percent)
	runtime.EventsEmit(pw.ctx, "percent", map[string]interface{}{
		"fileUUID": pw.ctx.Value("fileUUID"),
		"percent":  percent,
	})
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
