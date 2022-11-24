package web

import (
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileUpload struct {
	FileField string
	// 为什么要用户传
	// 要考虑文件重名
	DstPathFunc func(*multipart.FileHeader) string
}

func (u FileUpload) Handle() HandleFunc {

	if u.FileField == "" {
		u.FileField = "file"
	}
	if u.DstPathFunc == nil {
		// 设置默认值
		u.DstPathFunc = func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", uuid.New().String())
		}
	}

	return func(ctx *Context) {
		// 上传文件的逻辑代码在这里

		// 第一步: 读到文件内容
		// 第二部: 计算出目标路径
		// 第三步: 保存文件
		// 第四步: 返回响应
		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer file.Close()
		// 我怎么知道目标路径
		// 将目标路径计算的逻辑 交给用户
		dst := u.DstPathFunc(fileHeader)

		// 可以尝试把 dst上不存在的目录全部建立起来
		//os.MkdirAll()

		// os.O_WRONLY 写入数据
		// os.O_TRUNC 如果文件本身存在，那就清空数据
		// os.O_CREATE 如果文件不存在那就创建出来
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		defer dstFile.Close()
		// buffer 会影响你的性能
		// 要考虑复用
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败" + err.Error())
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}

// option 模式
type FileUploadOption func(unload *FileUpload)

func NewFileUploadOption(opts ...FileUploadOption) *FileUpload {
	res := &FileUpload{
		FileField: "file",
		DstPathFunc: func(header *multipart.FileHeader) string {
			return filepath.Join("testdata", "upload", uuid.New().String())
		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

type FileDownloader struct {
	// 下载的目标的路径
	Dir string
}

func (d FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		// 用的是 xxx?file=xxx
		req, err := ctx.QueryValue("file")
		if err != nil {
			ctx.RespStatusCode = http.StatusBadRequest
			ctx.RespData = []byte("找不到目标文件")
			return
		}
		req = filepath.Clean(req)
		dst := filepath.Join(d.Dir, req)
		fn := filepath.Base(dst)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")

		http.ServeFile(ctx.Resp, ctx.Req, dst)
	}
}
