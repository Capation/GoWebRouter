package web

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type FileUpload struct {
	FileField   string
	DstPathFunc func(*multipart.FileHeader) string
}

func (u FileUpload) Handle() HandleFunc {
	return func(ctx *Context) {
		// 上传文件的逻辑代码在这里

		// 第一步: 读到文件内容
		// 第二部: 计算出目标路径
		// 第三步: 保存文件
		// 第四步: 返回响应
		file, fileHeader, err := ctx.Req.FormFile(u.FileField)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			return
		}
		defer file.Close()
		// 我怎么知道目标路径
		// 将目标路径计算的逻辑 交给用户
		dst := u.DstPathFunc(fileHeader)
		// os.O_WRONLY 写入数据
		// os.O_TRUNC 如果文件本身存在，那就清空数据
		// os.O_CREATE 如果文件不存在那就创建出来
		dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			return
		}
		defer dstFile.Close()
		// buffer 会影响你的性能
		// 要考虑复用
		_, err = io.CopyBuffer(dstFile, file, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上传成功")
	}
}
