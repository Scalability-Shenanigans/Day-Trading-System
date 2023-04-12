import { StockHolding } from "../requests/requests";

export function canSellStock(
  userStockHoldings: StockHolding[],
  stock: string,
  amount: number,
  currentPrice: number
): boolean {
  for (const value of userStockHoldings) {
    if (value.Stock === stock) {
      if (value.Amount >= Math.round(amount / currentPrice)) {
        return true;
      } else {
        return false;
      }
    }
  }

  return false;
}
