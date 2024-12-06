# First Login Promotion
Develop a promotional campaign system for the Trinity app, enabling a 30% discount on Silver subscription plans for the first 100 users registering via campaign links. The system will generate time-limited vouchers to ensure efficient campaign management and user engagement.

## Technical decisions
1. Ngôn ngữ lập trình: Golang
 - Hiệu năng cao cho phát triển microservice
 - Sử dụng Goroutine để thực thi độc lập và dồng thời cho các tác vụ cần thiết. 
2. Web framework: Gin-Gonic
- Dễ dàng tiếp cận và phát triển Restful API
3. Database: PostgreSQL
- Lưu trữ dữ liệu và ràng buộc mối quan hệ của các bảng
4. Caching: Redis
- Tăng hiệu năng hệ thống thông qua caching.
- Giảm tải cho cơ sở dữ liệu với dữ liệu truy vấn thường xuyên
- Dùng làm nơi lưu trữ queue
5. API Design: Restful API
- Chuẩn hóa API giúp khả năng tăng mở rộng.
- Dùng Swagger để tạo tài liệu API ( API Documentation).
6. Deployment
- Sử dụng docker để đóng gói các service, đảm bảo tính nhất quán giữa môi trường development, staging và production...
- Sử dụng k8s hoặc docker compose để triển khai và quản lý các service
7. Security
- Input validation: Kiểm tra dữ liệu đầu vào để tránh các tấn công injection
- Rate limit: Giới hạn số lượng request.

## Assumptions made
1. Campaign logic:
 - Quy định số lượng người dùng được nhận ưu đãi tối đa khi đăng ký chiến dịch.
 - Đảm bảo mỗi người chỉ được nhận một voucher cho mỗi chiến dịch.
 - Voucher có thời hạn sử dụng cụ thể ( Ví dụ: thời hạn 7 ngày kể từ ngày tạo voucher. Hoặc thời gian kết thúc chiến dịch.)
2. Voucher system:
- Định dạng voucher code: mỗi voucher code sẽ là một chuỗi kí tự ngẫu nhiên sử dụng nanoid.
- Mỗi voucher là duy nhất và không thể tái sử dụng sau khi đã áp dụng.
- Mỗi voucher sẽ được áp dụng vào một sản phẩm nhất định để giảm giá sản phẩm theo ưu đãi của voucher.
- Voucher sẽ được xác thực thông qua Voucher Service để đảm bảo tính hợp lệ ( còn hạn, chưa sử dụng, chiến dịch...)
- Nếu voucher được hai người dùng cùng đăng ký trong cùng một khoảng thời gian, hệ thống đảm bảo chỉ tạo voucher cho đúng số lượt người đăng ký theo chiến dịch ( Ví dụ: giới hạn người dùng đăng ký cho chiến dịch là 100 người )
2. Tương tác với client:
- Người dùng sẽ được yêu cầu email để nhận ưu đãi theo chiến dịch thông qua API ``POST /api/v1/promo/register`` 
- Kiểm tra trạng thái voucher thông qua API ``GET /api/v1/voucher/{id}``

## Future improvements
Hệ thống hiện tại đã đáp ứng yêu cầu cơ bản, nhưng vẫn cần đề xuất một số cải thiện để nâng cao hiệu năng, tính mở rộng và khả năng bảo trì.
- Sử dụng k8s để tăng tốc độ truy vấn cũng như hiệu năng của dịch vụ.
- Xây dựng thêm API Gateway để kết nối đến nhiều microserive
- Chia nhỏ hệ thống thành nhiều microservice ( Ví dụ: User microservice, Product microservice...)
- Sử dụng Kafka để quản lý Queue cho các tác vụ như nhận voucher, thanh toán...
- Cho phép tạo hàng loạt voucher trước khi bắt đầu chiến dịch để giảm tải xử lý đồng thời.
- Real-time: Cung cấp dashboard hiển thị số lượng người dùng tham gia, tỷ lệ chuyển đổi, và tình trạng voucher theo thời gian thực.
- Security: Phân quyền truy cập API theo vai trò (RBAC)
- Có thể chuyển đổi thành GRPC để tăng khả năng bảo mật và đơn giản hóa giao tiếp truy vấn giữa các microservice.
- Hỗ trợ đa ngôn ngữ cho API.
- Thêm trương trình khuyến mãi nếu người dùng mời thêm bạn bè đăng ký và đã mua sản phẩm thành công.
- Tích hợp chức năng thanh toán online tử third-party, ví điện tử, visa card ...
- Sau khi thanh toán xong sẽ cập nhật trạng thái thanh toán vào bảng Order thông qua queue. Thông báo cho người dùng và quản trị viên biết trạng thái đơn hàng.

## Hướng dẫn cài đặt
1. Ngôn ngữ lập trình: Golang 1.21
2. Clone Repository:
```bash
  git clone git@github.com:nhanthanh93-com/first-login-promotion-test.git
  cd first-login-promotion-test
```

3. Tạo .env file từ thư mục gốc của dự án với nội dung như sau:
```.dotenv
POSTGRES_USER=dbadm
POSTGRES_PASSWORD=dbpwd
POSTGRES_DB=db

PORT=9888
APP_CONFIG=ewogICAgImRiQWRkciI6ICJob3N0PXRhX3Bvc3RncmVzIHVzZXI9ZGJhZG0gcGFzc3dvcmQ9ZGJwd2QgZGJuYW1lPWRiIHBvcnQ9NTQzMiBzc2xtb2RlPWRpc2FibGUgVGltZVpvbmU9QXNpYS9Ib19DaGlfTWluaCIsCiAgICAicmVkaXNBZGRyIjogInRhX3JlZGlzOjYzNzkiCn0=
```
Giải thích:
- POSTGRES_USER là tên user của database.
- POSTGRES_PASSWORD là mật khẩu của database.
- POSTGRES_DB là tên của database.
- APP_CONFIG là một chuỗi base64 string đã được endcode.
```json 
{
    "dbAddr": "host=ta_postgres user=dbadm password=dbpwd dbname=db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
    "redisAddr": "ta_redis:6379"
}
```
4. Cài đặt Dependencies
```bash
 go mod tidy 
```

5. Chạy service ở môi trường development 

- Sử dụng makeFile
```bash
    make dev
```

- Hoặc sử dụng docker
```bash
$ docker compose up
```

## API Documentation
Khởi tạo API Document bằng Swagger UI
- Sử dụng makeFile
```bash
    make swag_init
```
- Sử dụng swag cli
```bash
    swag init --parseDependency --parseInternal
```

- Truy cập URL sau để lấy API
```http://localhost:9888/swagger/doc.json```
- Truy cập URL sau để xem document bằng Swagger UI
```http://localhost:9888/swagger/index.html#/```

## Database Design
- dbdiagram.io
  ![db](https://github.com/user-attachments/assets/ce137f1e-7e2c-4228-ae8b-bed659ea76b9)

```text
Table campaign {
  id uuid [primary key]
  name varchar
  max_users integer
  created_at timestamp
  updated_at timestamp
  expires_at timestamp
}

Table voucher {
  id uuid [primary key]
  code varchar
  campaign_id uuid [ref: > campaign.id]
  user_id uuid [ref: > user.id]
  is_used boolean
  expires_at timestamp
  created_at timestamp
  updated_at timestamp
}

Table user {
  id uuid [primary key]
  email varchar
  created_at timestamp
  updated_at timestamp
}

Table cart {
  id uuid [primary key]
  user_id uuid [ref: > user.id]
  created_at timestamp
  updated_at timestamp
}

Table cart_item {
  id uuid [primary key]
  cart_id uuid [ref: > cart.id]
  product_id uuid [ref: > product.id]
  quantity int
  created_at timestamp
  updated_at timestamp
}

Table order {
  id uuid [primary key]
  user_id uuid [ref: > user.id]
  total float
  status string
  created_at timestamp
  updated_at timestamp
}


Table product {
  id uuid [primary key]
  name varchar
  price int
  created_at timestamp
  updated_at timestamp
}

Table campaign_product {
  id uuid [primary key]
  campaign_id uuid [ref: > campaign.id]
  product_id uuid [ref: > product.id]
}


```