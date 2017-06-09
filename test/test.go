package main
import(
    "fmt"
    "strings"
    "bytes"
)

func main(){
    //字符串拼接 strings.Join 接受[]string 可参考godoc
    v := []string{"a", "b", "c"}
    s := strings.Join(v, "-")
    fmt.Println(s)              // result: a-b-c

    //字符统计
    c := []byte("hello")
    count := bytes.Count(c, []byte("l"))
    fmt.Println(count)

}


