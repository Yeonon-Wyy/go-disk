package test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

func TestUppart(t *testing.T) {
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Inllb25vbiJ9.HsKMO_EznN-a9s-GsZABidMZ4hQNQyMc41a9jXtIt9ArH4BY3tON8Rno57EZFZVO41naXVmxZlkSkbwcr56tcnxyE8HwNY0ymYcgdzoY_z0nisRl1qwMtVatdeQGwQidy8MpMiYBRt5zjmoAIEuksNhLEXjTjY8iRTGnlEvEoSAYAePWP6yLKsPQm5yx7nesN979SK07JFItMCsqxsDy4Z6tHP3o7AUOFGoY3UrKjgK6NxrN6dXd_jp4ZyBqc5g2Qzh_Py4lf66cSE8cLpbs_DyAMcFwWHHYUcpPvjD0cN-eEHj6y9zCcqZNes-ccTHtJYOtbx-AVqT9huF_lpQoJw"
	tURL := "http://localhost:9000/files/mpupload/uppart?" +
		"username=yeonon"
	filename := "E:/333.zip"
	multipartUpload(filename, tURL, 5242880, token)

}

func multipartUpload(filename string, targetURL string, chunkSize int, token string) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++

		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)
		go func(b []byte, curIdx int) {
			client := &http.Client{}
			fmt.Printf("upload_size: %d\n", len(b))
			request, err := http.NewRequest("POST",
				targetURL + "&uploadid=yeonon15fcccfe6fd0de30" + "&index=" + strconv.Itoa(curIdx),
				bytes.NewReader(b),
			)

			if err != nil {
				fmt.Println(err)
				return
			}

			request.Header.Set("Content-Type", "multipart/form-data")
			request.Header.Set("Authorization", token)

			resp, err := client.Do(request)
			if err != nil {
				fmt.Println(err)
			}
			body, er := ioutil.ReadAll(resp.Body)
			fmt.Printf("%+v %+v\n", string(body), er)
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		//遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	for idx := 0; idx < index; idx++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		}
	}

	return nil
}