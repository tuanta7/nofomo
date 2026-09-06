# Đường Trung Bình Động

Đường Trung Bình Động (Moving Average - MA) là chỉ báo kỹ thuật phổ biến, được tính bằng giá trị trung bình của giá đóng cửa trong một số phiên nhất định.

- MA giúp làm mượt biến động giá, loại bỏ nhiễu ngắn hạn và xác định xu hướng thị trường.
- Có hai loại MA chính: MA Đơn Giản (SMA) và MA Hàm Mũ (EMA)

Giao cắt đường MA (MA Cross) là tín hiệu xảy ra khi đường MA ngắn hạn cắt đường MA dài hạn, thường được sử dụng để xác định điểm vào lệnh hoặc dự báo xu hướng.

- Golden Cross (Cắt lên - Tín hiệu mua): Khi MA ngắn hạn (ví dụ MA9) cắt lên MA dài hạn (ví dụ MA26), báo hiệu xu hướng tăng.
- Death Cross (Cắt xuống - Tín hiệu bán): Khi MA ngắn hạn cắt xuống MA dài hạn, cảnh báo đảo chiều giảm.

> [!NOTE]
> MA9 có nghĩa là tính giá trung bình 9 phiên gần nhất, còn MA26 tính trung bình 26 phiên.
>
> Độ dài phiên phụ thuộc vào khung thời gian chọn trên biểu dồ kỹ thuật, có thể là 1 ngày trên biểu đồ ngày,
> 1 giờ hoặc 4 giờ trên biểu đồ giờ, hoặc 5/15 phút trên biểu đồ phút

## Exponential Moving Average (EMA)

Tham khảo: [Đường EMA là gì?](https://www.vietcap.com.vn/kien-thuc/duong-ema-la-gi-ung-dung-cua-duong-ema-trong-giao-dich)

### EMA Crossover

Công thức chuẩn để tính EMA là

```txt
EMA hiện tại = (Giá đóng cửa hiện tại − EMA trước đó) × Hệ số làm mượt + EMA trước đó
```

- Giá đóng cửa hiện tại: là giá kết thúc của phiên giao dịch đang xét.
- EMA trước đó: là giá trị EMA đã tính được từ chu kỳ liền kề trước đó.
- Hệ số làm mượt (Smoothing Factor), ký hiệu là `α` hoặc `k` được tính theo công thức $\frac{2}{{N} + 1}$ với N là số chu kỳ EMA (9 ngày, 21 ngày,...)
