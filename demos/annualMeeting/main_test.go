package annualMeeting

import (
	"fmt"
	"github.com/kataras/iris/v12/httptest"
	"sync"
	"testing"
)

func TestMVC(t *testing.T){
	e := httptest.New(t, newApp())

	var wg sync.WaitGroup
	e.GET("/").Expect().Status(httptest.StatusOK).
		Body().Equal("当前总共参与抽奖的用户数为： 0")


	for i:=0;i<100;i++{
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.POST("/import").WithFormField("users",
				fmt.Sprintf("test_user%d",i)).Expect().
				Status(httptest.StatusOK)
		}(i)
	}
	wg.Wait()
	e.GET("/").Expect().Status(httptest.StatusOK).
		Body().Equal("当前总共参与抽奖的用户数为： 100")
	e.GET("/lucky").Expect().Status(httptest.StatusOK).
		Body().Equal("当前共有99个用户参与抽奖\n")
}
