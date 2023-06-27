# Golang’de Configuration Management

Herkese selamlar. Bu yazıda Golang ile configuration management ve 
environment konularından değineceğim. Bu yazıda, Viper’ı dosyadan 
veya environment variables ile nasıl kullanacağımızı öğreneceğiz.
Bir backend uygulaması geliştirirken genellikle development, 
testing, staging, ve production gibi farklı ortamlar için farklı 
yapılandırmalar kullanmamız gerekebilir. Bu yapılandırmayı bir 
dosyadan okumak, development ve test ortamlarımız için varsayılan 
yapılandırmayı kolayca belirlememizi sağlar. Farklı ortamlar için 
ayrı dosyalar oluşturarak, her ortam için özelleştirilmiş 
yapılandırmaları kolayca yönetebilirsiniz.

Ayrıca güvenlik açısından da bizlere avantajlar sağlar. 
Diyelim ki birçok özelliğe sahip bir uygulamanız var ve 
uygulamanızı veritabanına erişim sağlamak için DBURL, DBNAME, 
USERNAME ve PASSWORD gibi tüm veritabanı bilgilerini yapılandırdınız. 
Tüm bu bilgileri kod içerisine girerseniz, yetkisiz kişiler de DB’ye 
erişebilir. Eğer git gibi versiyon kontrol sistemi kullanıyorsanız, 
kodu gönderdiğinizde DB detayları herkese açık hale gelecektir. 
Bu durumda hassas bilgileri (örneğin, API anahtarları, veritabanı 
kimlik bilgileri vb.) bu bilgilerin kaynak koduyla birlikte 
yayınlanmasını önler. Bu sayede, bu bilgilerin yetkisiz kişiler 
tarafından erişilebilir olmasının önüne geçilir.

Sürdürülebilirlik açısından da bizlere avantaj sağlar. 
Proje ekibindeki herkesin ortak bir yapılandırma dosyası 
kullanması sayesinde, ekibin üyeleri arasında yapılandırma 
ayarlarının uyumlu ve tutarlı olmasını sağlar. Ayrıca, bu dosyanın 
sürdürülebilirliği daha yüksektir, çünkü yapılacak değişiklikler 
projenin kaynak kodunu etkilemez.

## Peki neden Viper?

Viper, bu amaçla kullanılan çok popüler bir Golang paketidir.

- Viper, dosyadan değerleri bulabilir, yükleyebilir 
ve ayrıştırabilir.
- JSON, TOML, YAML, ENV veya INI gibi birçok dosya türünü destekler.
- Varsayılan değerleri ayarlayabilir veya geçersiz kılabilir.
- Yaptığımız herhangi bir yapılandırma değişikliğini dosyaya 
kaydetmek için viper’ı kullanabiliriz.
- Şifrelenmiş ve şifrelenmemiş değerler için çalışır.
- Değerleri environment variables üzerinden veya command-line 
flag üzerinden okuyabilir.
- Ayrıca, uzak sistemlerde ayarlarınızı saklamayı tercih ederseniz, 
viper’ı doğrudan veri okumak için kullanabilirsiniz.

## Viper’ı yüklemek için

Terminalde aşağıda verdiğim komutu çalıştırdığınızda viper 
paketini sizin için indirecektir.

```bash
go get github.com/spf13/viper
```

Komutun tamamlanması ardından projemizdeki go.mod dosyasında 
viper’ı dependency olarak görebilirsiniz.

## Config dosyası oluşturalım

Şimdi, geliştirme için değerleri depolamak için 
app.env adında yeni bir dosya oluşturacağım.

```dotenv
SERVER_ADDRESS='0.0.0.0:8080'
SECRET_KEY='secret'
DB_DRIVER='postgres'
DB_SOURCE='postgres://username:password@localhost:5432/dbname?sslmode=disable'
```
## Config dosyamızı kodumuzun içinde tanımlayalım

envConfigs adında bir struct oluşturdum. 
Bu struct bizim tüm ortam değişkenlerimizi depolayacak.

```go
type envConfigs struct {
 ServerAddress string `mapstructure:"SERVER_ADDRESS"`
 SecretKey string `mapstructure:"SECRET_KEY"`
 DBDriver string `mapstructure:"DB_DRIVER"`
 DBSource string `mapstructure:"DB_SOURCE"`
}
```

Ardından okuduğumuz değişkenleri saklamak
için bu struct üzerinden bir değişken oluşturalım.

```go
var EnvConfigs *envConfigs
```

Bu değişkenleri okumak ve EnvConfigs değişkenine kaydetmek
için loadEnvVariables adında bir fonksiyon oluşturalım.
Bu fonksiyon bize yukarıda oluşturduğumuz envConfigs ile
dönüş yapacak.

```go
func loadEnvVariables() (config *envConfigs) {

viper.AddConfigPath(".")

viper.SetConfigName("app")

viper.SetConfigType("env")

err := viper.ReadInConfig()
if err != nil {
log.Panic("Error reading config file! ", err)
}

err = viper.Unmarshal(&config)
if err != nil {
log.Panic("Unable to decode into struct! ", err)
}

return
}
```

viper.AddConfigPath(“.”) ile değişkenleri tuttuğumuz
dosyanın konumunu belirtiyoruz. Direkt olarak projenin içine koyduğum için ben burada “.” olarak verdim.

viper.SetConfigName(“app”) ile dosyamızın
adını veriyoruz.

viper.SetConfigType(“env”) ile dosyamızın uzantısını
ve türünü veriyoruz. Oluşturduğunuz dosya .env den
farklı bir uzantıya sahipse burada onu vermelisiniz.

viper.ReadInConfig() ile dosyamızı okuyoruz.

viper.Unmarshal(&config) ile okuduğumuz verileri
oluşturduğumuz struct a çeviriyoruz.

```go
func InitEnvConfigs() {
EnvConfigs = loadEnvVariables()
}
```

Son olarak InitEnvConfigs adında bir fonksiyon
oluşturarak yukarıda oluşturmuş olduğumuz değişkene
bu değerleri atadık.

```go
package config

import "log"
import "github.com/spf13/viper"

var EnvConfigs *envConfigs

type envConfigs struct {
ServerAddress string `mapstructure:"SERVER_ADDRESS"`
SecretKey     string `mapstructure:"SECRET_KEY"`
DBDriver      string `mapstructure:"DB_DRIVER"`
DBSource      string `mapstructure:"DB_SOURCE"`
}

func InitEnvConfigs() {
EnvConfigs = loadEnvVariables()
}

func loadEnvVariables() (config *envConfigs) {

viper.AddConfigPath(".")

viper.SetConfigName("app")

viper.SetConfigType("env")

err := viper.ReadInConfig()
if err != nil {
log.Panic("Error reading config file! ", err)
}

err = viper.Unmarshal(&config)
if err != nil {
log.Panic("Unable to decode into struct! ", err)
}

return
}
```

Şimdi oluşturduğumuz ortam değişkenlerini kullanabiliriz.

```go
package main

import (
"fmt"
"golang-env/config"
)

func main() {

config.InitEnvConfigs()

fmt.Println("SERVER ADDRESS : ", config.EnvConfigs.ServerAddress)
fmt.Println("SECRET KEY : ", config.EnvConfigs.SecretKey)
fmt.Println("DB DRIVER : ", config.EnvConfigs.DBDriver)
fmt.Println("DB SOURCE : ", config.EnvConfigs.DBSource)
}
```

Main fonksiyonumuzun başında InitEnvConfigs
fonksiyonunu çağırarak bu değerlerin uygulamamız
başlarken okunup atanmasını sağlamış olduk.

```bash
go run main.go
```
```bash
SERVER ADDRESS :  0.0.0.0:8080
SECRET KEY :  secret
DB DRIVER :  postgres
DB SOURCE :  postgres://username:password@localhost/dbname?sslmode=disable
```

main.go yu çalıştırdığımızda ise yukarıdaki sonuca ulaşmış olduk.

.env değilde YAML veya JSON gibi formatlarda dosya
oluşturursanız tek yapmanız gereken viper.SetConfigType(“env”) yi
değiştirmek. “yaml” veya “json” değeri vererek kullanabilirsiniz.


