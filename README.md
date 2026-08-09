# yt-clone-server

Backend API cho project yt-clone, được viết bằng Go, Gin, GORM và PostgreSQL.

## Yêu cầu

- Go `1.26.5` hoặc phiên bản tương thích với `go.mod`
- PostgreSQL đang chạy ở máy local
- Git (không bắt buộc nếu source code đã có sẵn)

Kiểm tra Go đã được cài đặt:

```powershell
go version
```

### 1. Tải dependencies

```powershell
go mod download
go mod verify
```

Có thể dùng `go mod tidy` nếu muốn Go kiểm tra và đồng bộ lại dependencies:

```powershell
go mod tidy
```

### 2. Tạo database PostgreSQL

Tạo database có tên mặc định là `yt_clone`:

```powershell
createdb -U postgres yt_clone
```

Nếu lệnh `createdb` chưa có trong PATH, hãy tạo database bằng pgAdmin hoặc chạy SQL:

```sql
CREATE DATABASE yt_clone;
```

### 3. Tạo file `.env`

Tạo file `.env` trong thư mục `server` với nội dung sau. Thay `your_postgres_password` bằng mật khẩu PostgreSQL của bạn:

```env
PORT=8080
ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=yt_clone
DB_SSL_MODE=disable
```

`DB_PASSWORD` là biến bắt buộc. Các biến còn lại sẽ dùng giá trị mặc định nếu không được khai báo.

### 4. Khởi động server

```powershell
go run ./cmd/api
```

Khi chạy thành công, terminal sẽ hiển thị server đang chạy trên port `8080` và kết nối PostgreSQL thành công.

## Chạy server ở những lần tiếp theo

Mỗi lần mở terminal mới, đi tới thư mục server rồi chạy:

```powershell
cd D:\Code\Nextjs\yt-clone\server
go run ./cmd/api
```

Không cần chạy lại `go mod download` hoặc tạo database nếu dependencies và database đã được thiết lập trước đó.

Dừng server bằng:

```text
Ctrl + C
```

Mặc định server chạy tại:

```text
http://localhost:8080
```

Có thể đổi port trong `.env`, ví dụ:

```env
PORT=3000
```

## Build thành file executable

Nếu muốn build trước rồi chạy file executable:

```powershell
go build -o bin\api.exe .\cmd\api
.\bin\api.exe
```

## Một số lỗi thường gặp

### `Required environment variable is missing: DB_PASSWORD`

Chưa tạo `.env` hoặc chưa khai báo `DB_PASSWORD`. Kiểm tra file `.env` nằm đúng trong thư mục `server`.

### `Failed to connect to the database`

Kiểm tra:

- PostgreSQL đã được khởi động chưa
- `DB_HOST`, `DB_PORT`, `DB_USER` và `DB_PASSWORD` có đúng không
- Database `yt_clone` đã tồn tại chưa
- PostgreSQL có đang lắng nghe ở port `5432` không

### `go is not recognized`

Go chưa được cài đặt hoặc thư mục cài Go chưa được thêm vào PATH. Cài Go từ trang chính thức rồi mở lại terminal hoặc VS Code.

## Cấu trúc entrypoint

Server được khởi động từ:

```text
cmd/api/main.go
```

File này tải cấu hình, kết nối PostgreSQL, thiết lập Gin router và chạy HTTP server.
