package Server

import (
	Lenny_Gin "github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"../Model"
	"strconv"
	"github.com/dgrijalva/jwt-go"
	"time"
	
	"errors"
)

var (
	addr string
	readTimeOut int
	writeTimeOut int
	maxHeaderBytes uint32
	
)

func SetServerProperty(addr string, readtimeout int, writetimeout int, maxheaderbytes uint32) {
	addr = addr
	readTimeOut = readtimeout
	writeTimeOut = writetimeout
	maxHeaderBytes = maxheaderbytes
}

func init() {
	
	addr = ":8080"
	readTimeOut = 10
	writeTimeOut = 10
	maxHeaderBytes = 1 << 20
	
	fmt.Println("File gin_server_encapsulation init.")
}

func Lenny_Server_Start() {
	setRouter_Run()
}

func setRouter_Run() {
	router := Lenny_Gin.Default()
	Lenny_Gin.SetMode(Lenny_Gin.DebugMode)
	router.GET("/user", routerHandle_Get_User)
	router.POST("/user/registration", routerHandle_Post_Registration)
	router.POST("/user/login", routerHandle_Post_Login)
	
	
	router.Run(":9090")
}

func parsingToken(content *Lenny_Gin.Context) error {
	tokenString := content.GetHeader("token")
	if tokenString == "" {
		content.JSON(http.StatusForbidden, nil)
		return errors.New("没有Token！")
	}
	fmt.Printf("token:%v\n", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _ , ok  := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method:%v\n", token.Header["alg"])
		}
		return []byte("Lenny001"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["user_id"], claims["username"], claims["expire"])
	}else {
		fmt.Println(err)
	}
	
	
	return nil
}

//获取用户信息
func routerHandle_Get_User(context *Lenny_Gin.Context) {
	err :=  parsingToken(context)
	if err != nil {
		return
	}
	id := context.Query("id")
	if id == "" {
		context.JSON(http.StatusOK, Model.MakeJsonError(Model.Error_ParamsNotCorrect, "参数不对！"))
	}else {
		i, _ := strconv.ParseInt(id, 10, 64)
		u, err := Model.Get_A_User(i)
		if err != nil {
			context.JSON(http.StatusOK, Model.MakeJsonError(Model.Error_UserNotExsit, "没有这个用户！"))
			return
		}
		context.JSON(http.StatusOK, u)
	}
}
//注册
func routerHandle_Post_Registration(context *Lenny_Gin.Context) {
	u := &Model.User{}
	if context.PostForm("username") == "" || context.PostForm("password") == "" {
		context.JSON(http.StatusOK, Model.MakeJsonError(Model.Error_ParamsNotCorrect, "参数不对！"))
		return
	}
	u.Username = context.PostForm("username")
	u.Password = context.PostForm("password")
	err := u.SaveToDatabase()
	if err != nil {
		context.JSON(http.StatusOK, Model.MakeJsonError(Model.Error_UserAlreadyExsit, "用户已经存在！"))
	}else {
		context.JSON(http.StatusOK, Model.MakeJsonSuccess(true, "注册成功！"))
	}
}
//登录
func routerHandle_Post_Login(context *Lenny_Gin.Context) {
	u := &Model.User{}
	if context.PostForm("username") == "" || context.PostForm("password") == "" {
		context.JSON(http.StatusOK, Model.MakeJsonError(Model.Error_ParamsNotCorrect, "参数不对！"))
		return
	}
	u.Username = context.PostForm("username")
	//u.Password = context.PostForm("password")
	u = Model.Get_A_User_Info(u)
	if u == nil {
		context.JSON(http.StatusOK, Model.MakeJsonError(Model.Error_UserNotExsit, "用户不存在啊！"))
		return
	}
	fmt.Printf("uuuuuuuu%v\n", u)
	//用户存在， 查询Token表， 如果没有就生成Token
	t := &Model.Token{Username:u.Username, UserId:u.Id}
	t = Model.Get_A_UserToken_Info(t)
	fmt.Printf("tttttt%v\n", t)
	if t == nil {
		//没有Token ， 生成Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":u.Id,
			"username":u.Username,
			"expire": time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, err := token.SignedString([]byte("Lenny001"))
		fmt.Printf("%v\n tokenS:%v\n", err, tokenString)
		t = &Model.Token{UserId:u.Id,Username:u.Username,TokenString:tokenString}
		err = t.SaveToDatabase()
		if err != nil {
			context.JSON(http.StatusOK, Model.MakeJsonError(Model.Error_OperationFailed, "生成Token失败！"))
			return
		}
		//插入Token数据成功
		context.JSON(http.StatusOK, Model.Get_A_UserToken_Info(t))
	}else {
		//有Token， 更新Token的expire
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":u.Id,
			"username":u.Username,
			"expire": time.Now().Add(time.Hour * 24).Unix(),
		})
		tokenString, _ := token.SignedString([]byte("Lenny001"))
		t = Model.Get_A_UserToken_Info(t)
		t.TokenString = tokenString
		t.Update()
		//l := &Model.LoginSuccessToken{}
		
		context.JSON(http.StatusOK, Model.MakeJsonForLoginSuccess(true, "登录成功！", *t))
		
		
	}
}


