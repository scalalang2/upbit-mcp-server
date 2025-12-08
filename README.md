# upbit-mcp-server
내가 쓰기 위해서 만든 업비트 API MCP 서버

## 지원 기능
- 계정 관련 도구
  - `GetAccounts`: 전체 계좌 조회
  - `PlaceBuyOrder`: 시장가/지정가 매수 주문하기
  - `PlaceSellOrder`: 시장가/지정가 매도 주문하기
  - `CancelOrder`: 주문 취소
  - `GetAvailableOrderInfo`: 마켓 단위로 주문 가능 정보 확인
  - `GetClosedOrderHistory`: 완료된 주문 조회
  - `GetOpenOrderList`: 현재 진행중인 주문 리스트
- 시장 데이터 조회
  - `GetMarketSummary`: 특정 시장 정보 조회
  - `GetMarketTrends`: 현재 시장 트렌드 정보 조회 
  - `GetDayCandles`: 일봉
  - `GetWeekCandles`: 주봉
  - `GetMonthCandles`: 월봉
  - `GetMinuteCandles`: 분봉
- 기술적 분석 지표
  - 이동 평균선
  - MACD
  - RSI
  - Bollinger Bands
  - ODV

## MCP 연동 방법
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