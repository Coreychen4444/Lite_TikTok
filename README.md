# Lite_TikTok
A lite version TikTok and some basic functions included

## 视频及封面存储说明
由于本项目用Google Cloud Storage作为视频及封面存储，所以需要在Google Cloud Platform上创建一个项目，并创建一个存储桶，将main函数中的存储桶名改为自己的存储桶名
```go
bucketName := "your bucket name"
```
其中`"your bucket name"`为存储桶名
并且需要将Google Cloud Platform上的服务账号的密钥下载到本地，将json密钥文件放入项目根目录或设置全局变量，并将**service**目录下的**videoservice**文件中的`PublishVideo`函数中的密钥路径改为自己的密钥路径
```go
client, err := storage.NewClient(ctx, option.WithCredentialsFile())
	if err != nil {
		return fmt.Errorf("创建存储客户端失败")
	}
	defer client.Close()
```
其中`option.WithCredentialsFile()`中的参数为密钥路径，如果将密钥文件放入项目根目录，则参数为`"./*************.json"`，如果设置了全局变量，则参数为`os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")`。
可通过在终端shell中添加以下行来设置环境变量
```shell
export GOOGLE_APPLICATION_CREDENTIALS="your path/*************.json"
```
其中`"your path/*************.json"`为你存放在本地的密钥路径

## 数据库配置
确保已安装并启动mysql服务器，且已创建数据库tiktok，并且将main函数中的数据库配置改为自己的配置（用户名和密码）
```go
dsn := "root:44447777@tcp(127.0.0.1:3306)/tiktok_db?charset=utf8mb4&parseTime=True&loc=Local"
```
其中root为用户名，44447777为密码，tiktok_db为数据库名

## 运行
在项目根目录下终端运行
```shell
go run main.go
```
即可运行

## 效果展示
用安卓手机或模拟器安装tiktok简易版的apk文件，打开后双击“我”并输入服务器地址（如http://192.168.1.1:8080），即可使用
