# Thuật ngữ kỹ thuật

## 1. Chỉ báo kỹ thuật

Chỉ báo kỹ thuật (Indicators) là công thức toán học được áp dụng lên giá/volume để tạo ra tín hiệu giao dịch

Tick: Bước giá tối thiểu, chiều hướng thay đổi giá hoặc từng sự kiện giao dịch đơn lẻ

- Tick Size: Bước giá tối thiểu, khoảng cách tăng hoặc giảm giá nhỏ nhất mà một loại chứng khoán được phép dịch chuyển trên sàn giao dịch
- Tick Direction: Xác định giá dịch chuyển lên (Uptick) hay xuống (Downtick) so với giao dịch trước đó.
- Tick Data: Bản ghi chi tiết từng giao dịch khớp đơn lẻ trên thị trường
- Tick Volume: Chỉ báo đo lường số lần giá tài sản thay đổi trong một khoảng thời gian nhất định

Nến (Candle)

- Thông tin lịch sử nến (OHLCV - Open, High, Low, Close, Volume), được tạo từ một loạt tickdata
- Cây nến (Candlestick) trong biểu đồ là thể hiện của OHLCV (thường là trong 1 giây, 1 phút, 1 giờ, 1 ngày,...)

```stxt
OHLC DATA                          CANDLESTICK CHART
(Số liệu thô)                      (Hình ảnh trực quan)
        ↓                                  ↓
{                                  ┌───────────┐
  open: 25500,         ───────>    │   High    │ = 25525 <- Giá cao nhất
  high: 25525,                     │           │
  low: 25480,                      │  Body     │ = Open -> Close
  close: 25510,                    │(Red/Green)|   - Green: Close > Open ↑
  volume: 1500000                  │           │   - Red: Open > Close ↓
                                   │           │  
}                                  │  Low      │ = 25480 <- Giá thấp nhất
                                   └───────────┘
```

## 2. Danh mục đầu tư 

Danh mục đầu tư (Investment Portfolio) là tập hợp các tài sản tài chính mà một cá nhân hoặc tổ chức sở hữu

