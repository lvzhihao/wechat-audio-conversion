// Copyright © 2017 edwin <edwin.lzh@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/lvzhihao/goutils"
	"github.com/lvzhihao/wechat-audio-conversion/wxaudio"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "api server",
	Long:  `wxaudio api server`,
	Run: func(cmd *cobra.Command, args []string) {
		app := goutils.NewEcho()

		var logger *zap.Logger
		// app logger level
		if os.Getenv("DEBUG") == "true" {
			logger, _ = zap.NewDevelopment()
			app.Logger.SetLevel(log.DEBUG)
			app.Renderer = goutils.NewEchoRenderer("demo", "demo/views/*.html")
			app.GET("/demo", func(ctx echo.Context) error {
				return ctx.Render(http.StatusOK, "upload.html", map[string]string{
					"host":      viper.GetString("api_host"),
					"demo_code": viper.GetString("demo_code"),
				})
			})
		} else {
			logger, _ = zap.NewProduction()
		}

		store := wxaudio.NewStorage()
		store.Set(wxaudio.NewQiniuInstance(
			viper.GetString("qiniu_access_key"),
			viper.GetString("qiniu_secret_key"),
			viper.GetString("qiniu_zhiyakf_bucket"),
			viper.GetString("qiniu_zhiyakf_domain"),
			&storage.Config{
				// 空间对应的机房
				Zone: &storage.ZoneHuadong,
				// 是否使用https域名
				UseHTTPS: true,
				// 上传是否使用CDN上传加速
				UseCdnDomains: false,
			},
		))
		err := store.Link()
		if err != nil {
			logger.Fatal("New Storage Error", zap.Error(err))
		}

		tempdir, err := ioutil.TempDir("/tmp", viper.GetString("tempdir"))
		if err != nil {
			logger.Fatal("tempdir init error", zap.Error(err))
		}
		defer os.RemoveAll(tempdir) // clean up
		logger.Info("tempdir init success", zap.String("tempdir", tempdir))

		app.GET("/api/decoder", func(ctx echo.Context) error {
			source := ctx.QueryParam("source")
			if source == "" {
				return ctx.String(http.StatusBadRequest, "no source param")
			}
			rsp, err := http.Get(source)
			logger.Debug("request source", zap.Any("response_status", rsp.Status), zap.Error(err))
			if rsp.StatusCode != http.StatusOK {
				ctx.Logger().Error(err)
				return ctx.String(http.StatusBadRequest, "source unavailable")
			}
			b, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				ctx.Logger().Error(err)
				return ctx.JSON(http.StatusBadRequest, UploadError(err))
			}
			sourceFp, err := ioutil.TempFile(tempdir, "source")
			if err != nil {
				ctx.Logger().Error(err)
				return ctx.JSON(http.StatusInternalServerError, UploadError(err))
			}
			defer os.Remove(sourceFp.Name())
			start := bytes.Index(b, []byte("#!SILK_V3"))
			_, err = sourceFp.Write(b[start:]) //小U给的链接要舍弃第一个字节
			if err != nil {
				ctx.Logger().Error(err)
				return ctx.JSON(http.StatusInternalServerError, UploadError(err))
			}
			cmd := exec.Command("sbin/decoder", sourceFp.Name(), sourceFp.Name()+".pcm")
			err = cmd.Run()
			if err != nil {
				ctx.Logger().Error(err)
				return ctx.JSON(http.StatusInternalServerError, UploadError(err))
			}
			defer os.Remove(sourceFp.Name() + ".pcm")
			cmd = exec.Command("/usr/bin/ffmpeg", "-y", "-f", "s16le", "-ar", "24000", "-ac", "1", "-i", sourceFp.Name()+".pcm", sourceFp.Name()+".mp3")
			err = cmd.Run()
			if err != nil {
				ctx.Logger().Error(err)
				return ctx.JSON(http.StatusInternalServerError, UploadError(err))
			}
			defer os.Remove(sourceFp.Name() + ".mp3")
			params := make(map[string]string, 0)
			params["x:source"] = source
			rb, _ := ioutil.ReadFile(sourceFp.Name() + ".mp3")
			ret, err := store.Upload(getUploadFilename(goutils.RandStr(20)+".mp3"), rb, params)
			if err != nil {
				ctx.Logger().Error(err)
				return ctx.JSON(http.StatusInternalServerError, UploadError(err))
			}
			logger.Debug("upload result", zap.Any("result", ret), zap.Error(err))
			return ctx.JSON(http.StatusOK, ApiUploadReturn{
				Status: "success",
				Data:   []interface{}{ret},
			})
		})

		goutils.EchoStartWithGracefulShutdown(app, viper.GetString("api_host"))
	},
}

type ApiUploadReturn struct {
	Status string        `json:"status"`
	Error  error         `json:"msg"`
	Data   []interface{} `json:"data"`
}

func UploadError(err error) *ApiUploadReturn {
	return &ApiUploadReturn{
		Status: "error",
		Error:  err,
	}
}

func getUploadFilename(name string) string {
	return fmt.Sprintf("wxaudio/%d-%s", time.Now().UnixNano(), name)
}

func init() {
	RootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
