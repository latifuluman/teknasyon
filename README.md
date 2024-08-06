## TL;DR

```sh
Servisleri lokalde çalıştırma

git clone git@github.com:latifuluman/teknasyon.git
cd teknasyon/project
make up_build

Servisleri Test Etme

### 1. `user_create` İsteğini Gönderme
1. Postman'da Collection'ı açın.
2. `user_create` isteğini bulun ve açın.
3. İsteği gönderin ve yanıtı kontrol edin. Yanıtta kullanıcının başarıyla oluşturulduğunu doğrulayın.

### 2. `user_login` İsteğini Tamamlama
1. Postman'da Collection'ı açın.
2. `user_login` isteğini bulun ve açın.
3. `username` ve `password` alanlarını belirlediğiniz kullanıcı bilgileri ile doldurun.
4. İsteği gönderin ve dönen JWT token'ını kopyalayın.

### 3. `transfer_money` İsteğine JWT Token Ekleme
1. Postman'da Collection'ı açın.
2. `transfer_money` isteğini bulun ve açın.
3. `Authorization` sekmesine gidin ve tür olarak `Bearer Token` seçin.
4. `Token` alanına `user_login` isteğinden aldığınız JWT token'ını yapıştırın.
5. Transfer işlemi için gerekli verileri doldurun ve isteği gönderin.

### 4. Mailhog Mail Sunucusunu Kontrol Etme
1. Tarayıcınızı açın.
2. URL çubuğuna `http://localhost:8025` yazın ve Enter'a basın.
3. Mailhog arayüzü açılacaktır. Burada gönderilen mailleri görebilirsiniz.

### 5. MongoDB'ye Bağlanma ve Kayıtları Görüntüleme
1. MongoDB client'ınızı açın.
2. Aşağıdaki bağlantı stringini kullanarak MongoDB'ye bağlanın:
    ```
    mongodb://admin:password@localhost:27017/?authSource=admin
    ```
3. Bağlantı başarılı olduktan sonra veritabanındaki koleksiyonları ve kayıtları inceleyin.


```

![Mailhog Mails](./mailhog.png)
*Mailhog üzerinden kullanıcıya mail gönderme*


![Mongo DB Logs](./mongo.png)
*Transfer loglarının MongoDB de loglanması*


### For testing you can use the following command
```
go test ./...
```

# Mikroservis Mimari Tasarımı

## Kullanılan Servisler
- **account-service**
- **listener-service**
- **logger-service**
- **mail-service**
- **user-service**

## Kullanılan Teknolojiler
- **postgresql**
- **mongodb**
- **rabbitmq**
- **mailhog**
- **nginx**

## Mimarinin Genel Akışı
![Mikroservis Mimari Tasarımı](./system_design.png)

1. **Client**'tan gelen istekler **NGINX**'e yönlendirilir.
2. **NGINX**, gelen istekleri uygun olan servise yönlendirir ve load balancing işlemini gerçekleştirir.
3. İstekler ilgili mikroservisler tarafından işlenir ve gerekirse diğer servislerle iletişime geçilir.
4. Mikroservisler arası iletişim RabbitMQ ve gRPC kullanılarak sağlanır.
5. Loglama ve mail işlemleri için **logger-service** ve **mail-service** kullanılır.

## Servisler ve Detayları

### User Service
- **Port:** 80
- **Bağlantılar:** PostgreSQL
- **Endpointler:**
  - **POST** `/v1/users/login`: Kullanıcı bilgileriyle giriş yapar. Bilgiler doğruysa JWT token döner.
  - **POST** `/v1/users`: Kullanıcı ekler.
  - **GET** `/v1/users/{userID}`: İlgili kullnıcıyı döner.
  - **DELETE** `/v1/users/{userID}`: İlgili kullanıcıyı siler. 

### Account Service
- **Port:** 80
- **Bağlantılar:** PostgreSQL, RabbitMQ
- **İşlevler:**
  - Hesap işlemleriyle alakalı logları ve mail içeriğini RabbitMQ'ya gönderir.
- **Endpointler:**
  - **POST** `/v1/accounts`: Yeni bir hesap oluşturur.
  - **GET** `/v1/accounts`: Mevcut hesapları listeler.
  - **POST** `/v1/accounts/transfer/money`: İki hesap arasında para transferi yapar.
  - **GET** `/v1/accounts/{accountID}`: İlgili hesabı döner.

### Listener Service
- **Bağlantılar:** RabbitMQ
- **İşlevler:**
  - RabbitMQ'dan mail ve log eventlerini dinler.
  - Eğer mail türünde bir event varsa, 50001 portundan gRPC ile mail service'e gider.
  - Eğer log türünden bir event varsa, 50001 portundan gRPC ile logger service'e gider.

### Mail Service
- **Port:** 50001 (gRPC üzerinden dinler)
- **İşlevler:**
  - Listener service mail event'i gönderdiğinde, bu bilgilerle MailHog üzerinden mail atar.

### Logger Service
- **Port:** 50001 (gRPC üzerinden dinler)
- **İşlevler:**
  - Listener service log event'i gönderdiğinde, bu bilgileri MongoDB'ye log olarak kaydeder.

## Veritabanı ve Mesajlaşma Sistemleri
- **PostgreSQL:** User Service ve Account Service tarafından kullanılır.
- **MongoDB:** Logger Service tarafından logların kaydedilmesi için kullanılır.
- **RabbitMQ:** Eventlerin (mail ve log) kuyruğa alınması ve dinlenmesi için kullanılır.
- **MailHog:** Mail gönderimlerinin test edilmesi için kullanılır.
- **NGINX:** Load balancing ve isteklerin uygun servislere yönlendirilmesi için kullanılır.
