export interface Page<T> { items: T[]; total: number }
export interface DashboardSummary {
  transactions: number; clearingRecords: number; merchants: number;
  disputes: number; cards: number; tokens: number;
}
