package main

import(
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

type Todo struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Status bool `json:"status"`
}

//数据库连接
func initMySQL() (err error) {
	dsn := "root:123456@(localhost:3306)/db2?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}
	return DB.DB().Ping()
}

func main() {
	//连接数据库
	err := initMySQL()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	//绑定对应的表
	DB.AutoMigrate(&Todo{})

	//定义一个路由
	r := gin.Default()
	//告知路由静态文件的路径
	r.Static("/static", "static")
	//改制路由渲染模板的路径
	r.LoadHTMLGlob("templates/*")
	r.GET("/",func (c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	

	//v1
	v1Group := r.Group("v1")
	{
		//待办事项
		//添加
		v1Group.POST("/todo", func (c *gin.Context) {
			//前端获取数据
			var todo Todo
			c.BindJSON(&todo)
			//存入数据库
			if err = DB.Create(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//查看所有待办事项
		v1Group.GET("/todo", func (c *gin.Context) {
			var todoList []Todo
			if err = DB.Find(&todoList).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todoList)
			}
		})
		//修改某一个待办事项
		v1Group.PUT("/todo/:id", func (c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			var todo Todo
			if err = DB.Where("id=?",id).First(&todo).Error;err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
				return
			}
			c.BindJSON(&todo)
			if err = DB.Save(&todo).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, todo)
			}
		})
		//删除
		v1Group.DELETE("/todo/:id", func (c *gin.Context) {
			id, ok := c.Params.Get("id")
			if !ok {
				c.JSON(http.StatusOK, gin.H{"error": "无效的id"})
				return
			}
			if err = DB.Where("id=?", id).Delete(Todo{}).Error; err!= nil{
				c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{id:"deleted"})
			}
		})
	}
	r.Run()

	
}