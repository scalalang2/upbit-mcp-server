# upbit-mcp-server
MCP Server for Upbit APIs which is largest cryptocurency exchange in South Korea.

## Supported Functions
- `GetAccounts`: 전체 계좌 조회
- `GetOrderHistory`: 완료된 주문 조회
- `GetOrder`: 특정 주문 조회
- `GetOrders`: 주문 리스트 조회 (Open, Wait 등)
- `CancelOrder`: 주문 취소
- `PlaceOrder`: 주문하기
- `GetChance`: 주문 가능 정보 확인
- `GetCoinAddresses`: 전체 입금 주소 조회
- `GetCoinAddress`: 특정 코인 입금 주소 조회
- `GenerateCoinAddress`: 입금 주소 생성 요청
- `GetWithdraws`: 출금 리스트 조회
- `GetWithdraw`: 개별 출금 조회
- `GetWalletStatus`: 지갑 상태 조회
- `GetMarketCodes`: 마켓 코드 조회
- `GetTicks`: 최근 체결 내역
- `GetTicker`: 현재가 정보
- `GetOrderBooks`: 호가 정보
- `GetDayCandles`: 일봉
- `GetWeekCandles`: 주봉
- `GetMonthCandles`: 월봉
- `GetMinuteCandles`: 분봉

## How To Install
```json
{
  "mcpServers": {
    "Upbit": {
      "command": "Path to mcp server",
      "env": {
        "UPBIT_ACCESS_KEY": "Your upbit access key",
        "UPBIT_SECRET_KEY": "Your upbit secret key"
      }
    }
  }
}

```